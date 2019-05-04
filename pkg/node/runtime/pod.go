//
// KULADO INC. CONFIDENTIAL
// __________________
//
// [2014] - [2019] KULADO INC.
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of KULADO INC. and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to KULADO INC.
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from KULADO INC..
//

package runtime

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/node/envs"
	"github.com/onedomain/lastbackend/pkg/util/cleaner"
	"github.com/onedomain/lastbackend/pkg/util/filesystem"
)

const (
	logPodPrefix               = "node:runtime:pod:>"
	defaultRootLocalStorgePath = "/private/var/lib/lastbackend/runtime/"

	BUFFER_SIZE = 1024
)

// tplScript is a helper script this is added to the template the commands.
const logScript = `
echo ""
echo "[task: %s]"
echo ""
set -eux
%s
`

// logScript is a helper script that is added to
// the build script to logging run a command.
const tplScript = `
%s
`

func PodManage(ctx context.Context, key string, manifest *types.PodManifest) error {
	log.V(logLevel).Debugf("%s provision pod: %s", logPodPrefix, key)

	//==========================================================================
	// Destroy pod =============================================================
	//==========================================================================

	// Call destroy pod
	if manifest.State.Destroy {

		if task := envs.Get().GetState().Tasks().GetTask(key); task != nil {
			log.V(logLevel).Debugf("%s cancel pod creating: %s", logPodPrefix, key)
			task.Cancel()
		}

		p := envs.Get().GetState().Pods().GetPod(key)
		if p == nil {

			ps := types.NewPodStatus()
			ps.SetDestroyed()
			envs.Get().GetState().Pods().AddPod(key, ps)

			return nil
		}

		log.V(logLevel).Debugf("%s pod found > destroy it: %s", logPodPrefix, key)

		PodDestroy(ctx, key, p)

		p.SetDestroyed()
		envs.Get().GetState().Pods().SetPod(key, p)
		return nil
	}

	//==========================================================================
	// Check containers pod status =============================================
	//==========================================================================

	// Get pod list from current state
	p := envs.Get().GetState().Pods().GetPod(key)
	if p != nil {

		// restore pov volume claims
		podVolumeClaimRestore(key, manifest)

		switch true {
		case !PodSpecCheck(ctx, key, manifest) || len(manifest.Runtime.Tasks) > 0:
			PodDestroy(ctx, key, p)
			break
		case !PodVolumesCheck(ctx, key, manifest.Template.Volumes):
			log.Debugf("%s volumes data changed: %s", logPodPrefix, key)
			for _, v := range manifest.Template.Volumes {

				if v.Volume.Name != types.EmptyString {

					log.Debugf("%s attach volume %s for pod %s", logPodPrefix, v.Name, key)

					pv, err := PodVolumeAttach(ctx, key, v)
					if err != nil {
						log.Errorf("%s can not attach volume for pod: %s", logPodPrefix, err.Error())
						return err
					}

					p.Volumes[v.Name] = pv

				} else {

					log.Debugf("%s create pod volume %s for pod %s", logPodPrefix, v.Name, key)

					var name = podVolumeKeyCreate(key, v.Name)

					vol := envs.Get().GetState().Volumes().GetVolume(name)

					if vol == nil {
						log.V(logLevel).Debugf("%s update pod volume: volume not found: create %s: %s", logPodPrefix, key, v.Name)

						vs, err := PodVolumeCreate(ctx, key, v)
						if err != nil {
							log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
							return err
						}

						pv := &types.VolumeClaim{
							Name:   podVolumeClaimNameCreate(key, v.Name),
							Volume: name,
							Path:   vs.Status.Path,
						}

						envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)
						p.Volumes[pv.Name] = pv

					} else {

						_, err := PodVolumeUpdate(ctx, key, v)
						if err != nil {
							log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
							return err
						}
					}

				}

			}
			return PodRestart(ctx, key)
		default:
			return nil
		}
	}

	log.V(logLevel).Debugf("%s pod not found > create it: %s", logPodPrefix, key)

	ctx, cancel := context.WithCancel(context.Background())
	envs.Get().GetState().Tasks().AddTask(key, &types.NodeTask{Cancel: cancel})

	go func() {

		status, err := PodCreate(ctx, key, manifest)
		if err != nil {
			log.Errorf("%s can not create pod: %s err: %s", logPodPrefix, key, err.Error())
			status.SetError(err)
		}

		envs.Get().GetState().Pods().SetPod(key, status)
	}()

	return nil
}

func PodRestart(ctx context.Context, key string) error {

	pod := envs.Get().GetState().Pods().GetPod(key)
	if pod == nil {
		return errors.New("pod not found")
	}

	cri := envs.Get().GetCRI()

	for _, c := range pod.Runtime.Services {
		if err := cri.Restart(ctx, c.ID, nil); err != nil {
			return err
		}
	}

	return nil
}

func PodCreate(ctx context.Context, key string, manifest *types.PodManifest) (*types.PodStatus, error) {

	var (
		status    = types.NewPodStatus()
		namespace = getPodNamespace(key)

		setError = func(err error) (*types.PodStatus, error) {
			log.Errorf("%s can not pull image: %s", logPodPrefix, err)
			status.SetError(err)
			envs.Get().GetState().Pods().SetPod(key, status)
			PodClean(ctx, status)
			return status, err
		}
	)

	log.V(logLevel).Debugf("%s create pod: %s", logPodPrefix, key)

	//==========================================================================
	// Pull image ==============================================================
	//==========================================================================

	status.SetPull()

	envs.Get().GetState().Pods().AddPod(key, status)

	log.V(logLevel).Debugf("%s have %d volumes", logPodPrefix, len(manifest.Template.Volumes))
	for _, v := range manifest.Template.Volumes {

		var name string
		if v.Volume.Name != types.EmptyString {
			name = fmt.Sprintf("%s:%s", getPodNamespace(key), v.Volume.Name)
		} else {
			name = podVolumeKeyCreate(key, v.Name)
		}

		vol := envs.Get().GetState().Volumes().GetVolume(name)
		if vol == nil {
			log.V(logLevel).Debugf("%s update pod volume: volume not found: create %s: %s", logPodPrefix, key, v.Name)

			vs, err := PodVolumeCreate(ctx, key, v)
			if err != nil {
				log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
				return status, err
			}

			pv := &types.VolumeClaim{
				Name:   podVolumeClaimNameCreate(key, v.Name),
				Volume: name,
				Path:   vs.Status.Path,
			}

			envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)
			status.Volumes[pv.Name] = pv

		} else {

			_, err := PodVolumeUpdate(ctx, key, v)
			if err != nil {
				log.Errorf("%s can not update volume data: %s", logPodPrefix, err.Error())
				return status, err
			}

			claim := envs.Get().GetState().Volumes().GetClaim(podVolumeClaimNameCreate(key, v.Name))
			if claim == nil {
				pv := &types.VolumeClaim{
					Name:   podVolumeClaimNameCreate(key, v.Name),
					Volume: name,
					Path:   vol.Status.Path,
				}

				envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)
				status.Volumes[pv.Name] = pv
			}

		}

		envs.Get().GetState().Pods().SetPod(key, status)
	}

	if len(manifest.Runtime.Tasks) > 0 {
		for _, t := range manifest.Runtime.Tasks {
			pst := new(types.PodStatusPipelineStep)
			pst.Status = types.StateProvision
			pst.Error = false
			status.Runtime.Pipeline[t.Name] = pst
		}

		envs.Get().GetState().Pods().SetPod(key, status)
	}

	log.V(logLevel).Debugf("%s have %d images", logPodPrefix, len(manifest.Template.Containers))

	for _, c := range manifest.Template.Containers {
		log.V(logLevel).Debugf("%s pull image %s for pod if needed", logPodPrefix, c.Image.Name)
		if err := ImagePull(ctx, namespace, &c.Image); err != nil {
			log.Errorf("%s can not pull image: %s", logPodPrefix, err.Error())
			return setError(err)
		}
	}

	//==========================================================================
	// Run container ===========================================================
	//==========================================================================

	status.SetStarting()
	status.Steps[types.StepPull] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	envs.Get().GetState().Pods().SetPod(key, status)

	for _, s := range manifest.Template.Containers {
		for _, e := range s.EnvVars {
			if e.Secret.Name != types.EmptyString {
				log.V(logLevel).Debugf("%s get secret info from api", logPodPrefix)
				if err := SecretCreate(ctx, namespace, e.Secret.Name); err != nil {
					log.Errorf("%s can not fetch secret from api", logPodPrefix)
				}
			}
		}
	}

	var (
		primary  string
		services = make([]*types.ContainerManifest, 0)
	)

	if len(manifest.Runtime.Services) == 0 {
		for _, s := range manifest.Template.Containers {
			m, err := containerManifestCreate(ctx, key, s)
			if err != nil {
				log.Errorf("%s can not create container manifest from spec: %s", logPodPrefix, err.Error())
				return setError(err)
			}

			services = append(services, m)
		}
	}

	if len(manifest.Runtime.Services) != 0 {

		for _, name := range manifest.Runtime.Services {
			for _, s := range manifest.Template.Containers {

				if s.Name != name {
					continue
				}

				m, err := containerManifestCreate(ctx, key, s)
				if err != nil {
					log.Errorf("%s can not create container manifest from spec: %s", logPodPrefix, err.Error())
					return setError(err)
				}

				services = append(services, m)
			}
		}

	}

	// run services
	for _, svc := range services {

		if primary != types.EmptyString {
			svc.Network.Mode = fmt.Sprintf("container:%s", primary)
		} else {
			primary = svc.Name
		}

		if err := serviceStart(ctx, key, svc, status); err != nil {
			log.Errorf("%s can not start service: %s", logPodPrefix, err.Error())
			return status, err
		}

	}

	status.SetRunning()
	status.Steps[types.StepReady] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now().UTC(),
	}

	envs.Get().GetState().Pods().SetPod(key, status)

	// run tasks
	for _, t := range manifest.Runtime.Tasks {

		log.V(logLevel).Debugf("%s start task %s", logPodPrefix, t.Name)

		var f, e bool

		for _, s := range manifest.Template.Containers {

			if s.Name != t.Container {
				continue
			}

			f = true
			spec := *s

			if len(t.EnvVars) > 0 {
				for _, te := range t.EnvVars {
					var f = false
					for _, se := range spec.EnvVars {
						if te.Name == se.Name {
							se.Value = te.Value
							se.Secret = te.Secret
							se.Config = te.Config
							f = true
						}
					}
					if !f {
						spec.EnvVars = append(spec.EnvVars, te)
					}
				}
			}

			var buf bytes.Buffer
			for _, command := range t.Commands {
				buf.WriteString(fmt.Sprintf(tplScript, command))
			}

			escaped := fmt.Sprintf("%q", t.Name)
			escaped = strings.Replace(escaped, "$", `\$`, -1)
			script := fmt.Sprintf(logScript, escaped, buf.String())

			rootPath := defaultRootLocalStorgePath
			if len(envs.Get().GetConfig().Workdir) != 0 {
				rootPath = envs.Get().GetConfig().Workdir
			}

			filepath := path.Join(rootPath, strings.Replace(key, ":", "-", -1), "init")

			log.Debugf("%s runtime volume create: %s", logPodPrefix, filepath)

			err := podLocalFileCreate(filepath, script)
			if err != nil {
				log.Errorf("%s can not create runtime volume err: %s", logPodPrefix, err.Error())
				return setError(err)
			}

			log.Debugf("%s container manifest create", logPodPrefix)

			m, err := containerManifestCreate(ctx, key, &spec)
			if err != nil {
				log.Errorf("%s can not create container manifest from spec: %s", logPodPrefix, err.Error())
				return setError(err)
			}

			if primary != types.EmptyString {
				m.Network.Mode = fmt.Sprintf("container:%s", primary)
			}

			m.Name = ""

			m.Exec.Command = []string{"/usr/local/bin/lb_entrypoint"}
			m.Exec.Entrypoint = []string{"/bin/sh"}
			m.Binds = append(m.Binds, fmt.Sprintf("%s:%s:ro", filepath, "/usr/local/bin/lb_entrypoint"))

			if err := taskExecute(ctx, key, t, *m, status); err != nil {
				log.Errorf("%s can not execute task: %s", logPodPrefix, err.Error())
				return setError(err)
			}

			for _, s := range status.Runtime.Pipeline {
				if s.Error {
					e = true
					status.SetError(errors.New(s.Message))
					break
				}
			}

		}

		if e || !f {
			break
		}
	}

	if len(manifest.Runtime.Tasks) > 0 {
		PodExit(ctx, key, status, true)
	}

	return status, nil
}

func PodClean(ctx context.Context, status *types.PodStatus) {

	for _, c := range status.Runtime.Services {
		log.V(logLevel).Debugf("%s remove unnecessary container: %s", logPodPrefix, c.ID)
		if err := envs.Get().GetCRI().Remove(ctx, c.ID, true, true); err != nil {
			log.Warnf("%s can-not remove unnecessary container %s: %s", logPodPrefix, c.ID, err)
		}
	}

	for _, c := range status.Runtime.Services {
		log.V(logLevel).Debugf("%s try to clean image: %s", logPodPrefix, c.Image.Name)
		//if err := ImageRemove(ctx, c.Image.Name); err != nil {
		//	log.Errorf("%s can not remove image: %s", logPodPrefix, err.Error())
		//	continue
		//}
	}

}

func PodExit(ctx context.Context, pod string, status *types.PodStatus, clean bool) {
	log.V(logLevel).Debugf("%s exit pod: %s", logPodPrefix, pod)

	timeout := time.Duration(3 * time.Second)
	for _, c := range status.Runtime.Services {

		var attempts = 5
		for i := 1; i <= attempts; i++ {
			if err := envs.Get().GetCRI().Stop(ctx, c.ID, &timeout); err != nil {
				// TODO: check container not found error
				log.Warnf("%s can-not stop container %s: %s", logPodPrefix, c.ID, err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}

		c.State.Stopped = types.PodContainerStateStopped{
			Stopped: true,
			Exit: types.PodContainerStateExit{
				Code:      0,
				Timestamp: time.Now(),
			},
		}
	}

	status.Steps[types.StateExited] = types.PodStep{
		Ready:     true,
		Timestamp: time.Now(),
	}

	if status.Status != types.StateError {
		status.SetExited()
	}

	envs.Get().GetState().Pods().SetPod(pod, status)

	if clean {
		PodClean(ctx, status)
		return
	}
}

func PodDestroy(ctx context.Context, pod string, status *types.PodStatus) {
	log.V(logLevel).Debugf("%s try to remove pod: %s", logPodPrefix, pod)
	PodClean(ctx, status)
	envs.Get().GetState().Pods().DelPod(pod)
	for _, v := range status.Volumes {
		if err := PodVolumeDestroy(ctx, pod, v.Name); err != nil {
			log.Errorf("%s can not destroy pod: %s", logPodPrefix, err.Error())
		}
	}

	rootPath := defaultRootLocalStorgePath
	if len(envs.Get().GetConfig().Workdir) != 0 {
		rootPath = envs.Get().GetConfig().Workdir
	}

	dirPath := path.Join(rootPath, strings.Replace(pod, ":", "-", -1), "init")

	log.Debugf("%s runtime volume remove: %s", logPodPrefix, dirPath)

	if err := podLocalFileDestroy(dirPath); err != nil {
		log.Errorf("%s can not destroy runtime volume path: %s", logPodPrefix, err.Error())
	}
}

func PodRestore(ctx context.Context) error {

	log.V(logLevel).Debugf("%s runtime restore state", logPodPrefix)

	cl, err := envs.Get().GetCRI().List(ctx, true)
	if err != nil {
		log.Errorf("%s pods restore error: %s", logPodPrefix, err)
		return err
	}

	for _, c := range cl {

		if c.Pod == types.EmptyString {
			continue
		}

		log.V(logLevel).Debugf("%s pod [%s] > container restore %s", logPodPrefix, c.Pod, c.ID)

		status := envs.Get().GetState().Pods().GetPod(c.Pod)
		if status == nil {
			status = types.NewPodStatus()
		}

		key := c.Pod

		cs := &types.PodContainer{
			ID:   c.ID,
			Name: c.Name,
			Image: types.PodContainerImage{
				Name: c.Image,
			},
			Envs:  c.Envs,
			Ports: c.Network.Ports,
			Binds: c.Binds,
		}

		cs.Restart.Policy = c.Restart.Policy
		cs.Restart.Attempt = c.Restart.Retry
		cs.Exec = c.Exec

		switch c.State {
		case types.StateCreated:
			cs.State = types.PodContainerState{
				Created: types.PodContainerStateCreated{
					Created: time.Now().UTC(),
				},
			}
		case types.StateStarted:
			cs.State = types.PodContainerState{
				Started: types.PodContainerStateStarted{
					Started:   true,
					Timestamp: time.Now().UTC(),
				},
			}
			cs.State.Stopped.Stopped = false
		case types.StatusStopped:
			cs.State.Stopped.Stopped = true
			cs.State.Stopped.Exit = types.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
		case types.StateError:

			cs.State.Error.Error = true
			cs.State.Error.Message = c.Status
			cs.State.Error.Exit = types.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
			cs.State.Stopped.Stopped = false
			cs.State.Stopped.Exit = types.PodContainerStateExit{
				Code:      c.ExitCode,
				Timestamp: time.Now().UTC(),
			}
			cs.State.Started.Started = false
		}

		if c.Status == types.StatusStopped {
			cs.State.Stopped = types.PodContainerStateStopped{
				Stopped: true,
				Exit: types.PodContainerStateExit{
					Code:      c.ExitCode,
					Timestamp: time.Now().UTC(),
				},
			}
		}

		cs.Ready = true
		status.Runtime.Services[cs.ID] = cs
		status.Network.PodIP = c.Network.IPAddress

		log.V(logLevel).Debugf("%s container restored %s", logPodPrefix, c.ID)
		envs.Get().GetState().Pods().SetPod(key, status)
		log.V(logLevel).Debugf("%s Pod restored %s: %s", key, status.State)
	}

	return nil
}

func PodLogs(ctx context.Context, id string, follow bool, s io.Writer, doneChan chan bool) error {

	log.V(logLevel).Debugf("%s get container [%s] logs streaming", logPodPrefix, id)

	var (
		cri    = envs.Get().GetCRI()
		buffer = make([]byte, BUFFER_SIZE)
		done   = make(chan bool, 1)
	)

	req, err := cri.Logs(ctx, id, true, true, follow)
	if err != nil {
		log.Errorf("%s error get logs stream %s", logPodPrefix, err)
		return err
	}
	defer func() {
		log.V(logLevel).Debugf("%s stop container [%s] logs streaming", logPodPrefix, id)
		ctx.Done()
		close(done)
		req.Close()
	}()

	go func() {
		for {
			select {
			case <-done:
				req.Close()
				return
			default:

				n, err := cleaner.NewReader(req).Read(buffer)
				if err != nil {

					if err == context.Canceled {
						log.V(logLevel).Debugf("%s Stream is canceled", logPodPrefix)
						return
					}

					log.Errorf("%s read bytes from stream err %s", logPodPrefix, err)
					doneChan <- true
					return
				}

				_, err = func(p []byte) (n int, err error) {
					n, err = s.Write(p)
					if err != nil {
						log.Errorf("%s write bytes to stream err %s", logPodPrefix, err)
						return n, err
					}

					if f, ok := s.(http.Flusher); ok {
						f.Flush()
					}
					return n, nil
				}(buffer[0:n])

				if err != nil {
					log.Errorf("%s write to stream err: %s", logPodPrefix, err)
					done <- true
					return
				}

				for i := 0; i < n; i++ {
					buffer[i] = 0
				}
			}
		}
	}()

	<-doneChan

	return nil
}

func PodSpecCheck(ctx context.Context, key string, manifest *types.PodManifest) bool {

	log.V(logLevel).Infof("%s pod check spec pod: %s", logPodPrefix, key)

	state := envs.Get().GetState().Pods().GetPod(key)

	var statec = make(map[string]*types.ContainerManifest, 0)
	var specc = make(map[string]*types.ContainerManifest, 0)

	for _, c := range manifest.Template.Containers {
		mf, err := containerManifestCreate(ctx, key, c)
		if err != nil {
			return false
		}
		specc[mf.Name] = mf
	}

	for _, c := range state.Runtime.Services {
		statec[c.Name] = c.GetManifest()
	}

	if len(statec) != len(specc) {
		log.Debugf("%s container spec count not equal not exists: %d != %d", logPodPrefix, len(statec), len(specc))
		return false
	}

	for n, mf := range specc {

		if _, ok := statec[n]; !ok {
			log.Debugf("%s container spec not exists: %s", logPodPrefix, n)
			return false
		}

		// check image

		c := statec[n]

		if c.Image != mf.Image {
			log.Debugf("%s images not equal: %s != %s", logPodPrefix, c.Image, mf.Image)
			return false
		}

		img := envs.Get().GetState().Images().GetImage(c.Image)
		if img == nil {
			log.Debugf("%s image not found in state: %s", logPodPrefix, mf.Image)
			return false
		}

		if len(mf.Exec.Command) == 0 {
			if strings.Join(c.Exec.Command, " ") != strings.Join(img.Status.Container.Exec.Command, " ") {
				log.Debugf("%s cmd different with img cmd: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Command, " "),
					strings.Join(img.Status.Container.Exec.Command, " "))
				return false
			}
		} else {
			if strings.Join(c.Exec.Command, " ") != strings.Join(mf.Exec.Command, " ") {
				log.Debugf("%s cmd different with manifest cmd: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Command, " "),
					strings.Join(mf.Exec.Command, " "))
				return false
			}
		}

		if len(mf.Exec.Entrypoint) == 0 {
			if strings.Join(c.Exec.Entrypoint, " ") != strings.Join(img.Status.Container.Exec.Entrypoint, " ") {
				log.Debugf("%s entrypoint changed: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Entrypoint, " "),
					strings.Join(img.Status.Container.Exec.Entrypoint, " "))
				return false
			}
		} else {
			if strings.Join(c.Exec.Entrypoint, " ") != strings.Join(mf.Exec.Entrypoint, " ") {
				log.Debugf("%s entrypoint changed: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Entrypoint, " "),
					strings.Join(mf.Exec.Entrypoint, " "))
				return false
			}
		}

		if mf.Exec.Workdir == types.EmptyString {
			if c.Exec.Workdir != img.Status.Container.Exec.Workdir {
				log.Debugf("%s workdir changed: %s != %s", logPodPrefix, c.Exec.Workdir, img.Status.Container.Exec.Workdir)
				return false
			}
		} else {
			if c.Exec.Workdir != mf.Exec.Workdir {
				log.Debugf("%s workdir changed: %s != %s", logPodPrefix, c.Exec.Workdir, mf.Exec.Workdir)
				return false
			}
		}

		if len(mf.Exec.Args) != 0 {
			if strings.Join(c.Exec.Args, " ") != strings.Join(mf.Exec.Args, " ") {
				log.Debugf("%s args changed: %s != %s", logPodPrefix,
					strings.Join(c.Exec.Args, " "),
					strings.Join(mf.Exec.Args, " "))
				return false
			}
		}

		// Check environments
		for _, e := range mf.Envs {
			var f = false
			for _, ie := range c.Envs {

				if ie == e {
					f = true
					break
				}
			}
			if !f {
				log.Debugf("%s env not found: %s", logPodPrefix, e)
				return false
			}
		}

		for _, e := range c.Envs {
			var f = false
			for _, ie := range mf.Envs {
				if ie == e {
					f = true
					break
				}
			}

			if !f {
				for _, ie := range img.Status.Container.Envs {
					if ie == e {
						f = true
						break
					}
				}
			}

			if !f {
				log.Debugf("%s env is unnecessary: %s", logPodPrefix, e)
				return false
			}
		}

		// Check binds
		for _, e := range mf.Binds {
			var f = false
			for _, ie := range c.Binds {
				if ie == e {
					f = true
					break
				}
			}
			if !f {
				log.Debugf("%s bind not found: %s", logPodPrefix, e)
				return false
			}
		}

		for _, e := range c.Binds {
			var f = false
			for _, ie := range mf.Binds {
				if ie == e {
					f = true
					break
				}
			}
			if !f {
				log.Debugf("%s bind is unnecessary: %s", logPodPrefix, e)
				return false
			}
		}

		// Check ports
		for _, e := range mf.Ports {
			var f = false
			for _, ie := range c.Ports {
				if e.HostIP != types.EmptyString {
					if e.HostIP != ie.HostIP {
						log.Debugf("%s port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s", logPodPrefix,
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				}

				if e.Protocol != types.EmptyString {
					if e.Protocol != ie.Protocol {
						log.Debugf("%s port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s", logPodPrefix,
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				} else {
					if ie.Protocol != "tcp" {
						log.Debugf("%s port map check failed: \t\t %s:%d:%d/%s == %s:%d:%d/%s", logPodPrefix,
							e.HostIP, e.HostPort, e.ContainerPort, e.Protocol,
							ie.HostIP, ie.HostPort, ie.ContainerPort, ie.Protocol)
						return false
					}
				}

				if ie.ContainerPort == e.ContainerPort &&
					ie.HostPort == ie.HostPort {
					f = true
					break
				}
			}

			if !f {
				log.Debugf("%s port map not found: \t\t %s:%d:%d/%s ", logPodPrefix,
					e.HostIP, e.HostPort, e.ContainerPort, e.Protocol)
				return false
			}
		}

		for _, e := range c.Ports {
			var f = false
			for _, ie := range mf.Ports {
				if ie.ContainerPort == e.ContainerPort &&
					ie.HostPort == ie.HostPort &&
					ie.Protocol == ie.Protocol &&
					ie.HostIP == ie.HostIP {
					f = true
					break
				}
			}

			if !f {
				log.Debugf("%s port map is unnecessary: %d", logPodPrefix, e.HostPort)
				return false
			}
		}

		if mf.RestartPolicy.Policy != c.RestartPolicy.Policy ||
			mf.RestartPolicy.Attempt != c.RestartPolicy.Attempt {

			log.Debugf("%s restart policy changed: %s:%d => %s:%d", logPodPrefix,
				c.RestartPolicy.Policy, c.RestartPolicy.Attempt,
				mf.RestartPolicy.Policy, mf.RestartPolicy.Attempt)
			return false
		}

	}

	return true

}

func PodVolumesCheck(ctx context.Context, pod string, spec []*types.SpecTemplateVolume) bool {

	log.V(logLevel).Debugf("%s check pod volumes: %s: %d", logPodPrefix, pod, len(spec))

	for _, v := range spec {

		if v.Volume.Name != types.EmptyString {
			continue
		}

		name := podVolumeKeyCreate(pod, v.Name)

		if v.Config.Name != types.EmptyString && len(v.Config.Binds) > 0 {
			equal, err := VolumeCheckConfigData(ctx, name, v.Config.Name)
			if err != nil {
				return false
			}
			return equal
		}

		if v.Secret.Name != types.EmptyString && len(v.Secret.Binds) > 0 {
			equal, err := VolumeCheckSecretData(ctx, name, v.Config.Name)
			if err != nil {
				return false
			}
			return equal
		}
	}

	return true
}

func PodVolumeUpdate(ctx context.Context, pod string, spec *types.SpecTemplateVolume) (*types.VolumeStatus, error) {

	log.V(logLevel).Debugf("%s update pod volume: %s: %s", logPodPrefix, pod, spec.Name)

	path := strings.Replace(pod, ":", "-", -1)
	path = fmt.Sprintf("%s-%s", path, spec.Name)

	var (
		name = podVolumeKeyCreate(pod, spec.Name)
	)

	status := envs.Get().GetState().Volumes().GetVolume(name)

	if spec.Secret.Name != types.EmptyString && len(spec.Secret.Binds) > 0 {
		if err := VolumeSetSecretData(ctx, name, spec.Secret.Name); err != nil {
			log.Errorf("%s can not set config data to volume: %s", logPodPrefix, err.Error())
			return status, err
		}
	}

	if spec.Secret.Name == types.EmptyString && spec.Config.Name != types.EmptyString && len(spec.Config.Binds) > 0 {
		if err := VolumeSetConfigData(ctx, name, spec.Config.Name); err != nil {
			log.Errorf("%s can not set config data to volume: %s", logPodPrefix, err.Error())
			return status, err
		}
	}

	return status, nil
}

func PodVolumeAttach(ctx context.Context, pod string, spec *types.SpecTemplateVolume) (*types.VolumeClaim, error) {

	log.V(logLevel).Debugf("%s attach pod volume: %s: %s", logPodPrefix, pod, spec.Name)

	var name = fmt.Sprintf("%s:%s", getPodNamespace(pod), spec.Name)

	volume := envs.Get().GetState().Volumes().GetVolume(name)
	if volume == nil {
		return nil, errors.New("volume not found on node")
	}

	pv := &types.VolumeClaim{
		Name:   podVolumeClaimNameCreate(pod, spec.Name),
		Volume: name,
		Path:   volume.Status.Path,
	}

	envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)

	return pv, nil
}

func PodVolumeCreate(ctx context.Context, pod string, spec *types.SpecTemplateVolume) (*types.VolumeStatus, error) {

	log.V(logLevel).Debugf("%s create pod volume: %s:%s", logPodPrefix, pod, spec.Name)

	hostPath := strings.Replace(pod, ":", "-", -1)
	hostPath = fmt.Sprintf("%s-%s", hostPath, spec.Name)

	var (
		name = podVolumeKeyCreate(pod, spec.Name)
		vm   = types.VolumeManifest{
			HostPath: hostPath,
			Type:     types.KindVolumeHostDir,
		}
	)

	st, err := VolumeCreate(ctx, name, &vm)
	if err != nil {
		log.Errorf("%s can not create pod volume: %s", logPodPrefix, err.Error())
		return nil, err
	}

	if spec.Secret.Name != types.EmptyString && len(spec.Secret.Binds) > 0 {
		if err := VolumeSetSecretData(ctx, name, spec.Secret.Name); err != nil {
			log.Errorf("%s can not set secret data to volume: %s", logPodPrefix, err.Error())
			return st, err
		}
	}

	if spec.Secret.Name == types.EmptyString && spec.Config.Name != types.EmptyString && len(spec.Config.Binds) > 0 {
		if err := VolumeSetConfigData(ctx, name, spec.Config.Name); err != nil {
			log.Errorf("%s can not set config data to volume: %s", logPodPrefix, err.Error())
			return st, err
		}
	}

	envs.Get().GetState().Volumes().SetLocal(name)

	return st, nil
}

func PodVolumeDestroy(ctx context.Context, pod, volume string) error {
	envs.Get().GetState().Volumes().DelLocal(podVolumeKeyCreate(pod, volume))
	return VolumeDestroy(ctx, podVolumeKeyCreate(pod, volume))
}

func podVolumeClaimRestore(key string, manifest *types.PodManifest) {

	pod := envs.Get().GetState().Pods().GetPod(key)
	if pod == nil {
		return
	}

	for _, v := range manifest.Template.Volumes {

		var name string
		if v.Volume.Name != types.EmptyString {
			name = fmt.Sprintf("%s:%s", getPodNamespace(key), v.Volume.Name)
		} else {
			name = podVolumeKeyCreate(key, v.Name)
		}

		vol := envs.Get().GetState().Volumes().GetVolume(name)
		if vol == nil {
			continue
		}

		claim := envs.Get().GetState().Volumes().GetClaim(podVolumeClaimNameCreate(key, v.Name))
		if claim == nil {
			pv := &types.VolumeClaim{
				Name:   podVolumeClaimNameCreate(key, v.Name),
				Volume: name,
				Path:   vol.Status.Path,
			}

			envs.Get().GetState().Volumes().SetClaim(pv.Name, pv)
			pod.Volumes[pv.Name] = pv
		} else {
			pod.Volumes[claim.Name] = claim
		}
	}
}

func podVolumeKeyCreate(pod, volume string) string {
	return fmt.Sprintf("%s-%s", strings.Replace(pod, ":", "-", -1), volume)
}

func podVolumeClaimNameCreate(pod, volume string) string {
	return fmt.Sprintf("%s:%s", pod, volume)
}

func podLocalFileCreate(path, data string) error {
	return filesystem.WriteStrToFile(path, data, 0777)
}

func podLocalFileDestroy(path string) error {
	return os.RemoveAll(path)
}

// TODO: move to distribution
func getPodNamespace(key string) string {
	var namespace = types.DEFAULT_NAMESPACE

	parts := strings.Split(key, ":")

	if len(parts) == 4 {
		namespace = parts[0]
	}

	return namespace
}
