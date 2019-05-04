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

package types

import (
	"golang.org/x/oauth2"
	"time"
)

type Vendor struct {
	Host  string
	Name  string
	Token *oauth2.Token
}

type User struct {
	Username  string
	ServiceID string
}

type VCSRepository struct {
	Name          string
	Description   string
	Private       bool
	DefaultBranch string
	Permissions   struct {
		Admin bool
	}
}

type VCSRepositories []VCSRepository

type VCSBranch struct {
	Name       string
	LastCommit Commit
}

type VCSBranches []VCSBranch

type Commit struct {
	Username string
	Hash     string
	Message  string
	Date     time.Time
	Email    string
}

type Commits []Commit
