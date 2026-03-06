/*
 * MinIO Cloud Storage, (C) 2018 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nas

import (
	"context"

	obstor "github.com/cloudment/obstor/cmd"
	"github.com/cloudment/obstor/pkg/auth"
	"github.com/cloudment/obstor/pkg/madmin"
	"github.com/urfave/cli"
)

func init() {
	const nasGatewayTemplate = `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS]{{end}} PATH
{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
PATH:
  path to NAS mount point

EXAMPLES:
  1. Start obstor gateway server for NAS backend
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_ROOT_USER{{.AssignmentOperator}}accesskey
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_ROOT_PASSWORD{{.AssignmentOperator}}secretkey
     {{.Prompt}} {{.HelpName}} /shared/nasvol

  2. Start obstor gateway server for NAS with edge caching enabled
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_ROOT_USER{{.AssignmentOperator}}accesskey
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_ROOT_PASSWORD{{.AssignmentOperator}}secretkey
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_CACHE_DRIVES{{.AssignmentOperator}}"/mnt/drive1,/mnt/drive2,/mnt/drive3,/mnt/drive4"
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_CACHE_EXCLUDE{{.AssignmentOperator}}"bucket1/*,*.png"
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_CACHE_QUOTA{{.AssignmentOperator}}90
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_CACHE_AFTER{{.AssignmentOperator}}3
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_CACHE_WATERMARK_LOW{{.AssignmentOperator}}75
     {{.Prompt}} {{.EnvVarSetCommand}} OBSTOR_CACHE_WATERMARK_HIGH{{.AssignmentOperator}}85
     {{.Prompt}} {{.HelpName}} /shared/nasvol
`

	obstor.RegisterGatewayCommand(cli.Command{
		Name:               obstor.NASBackendGateway,
		Usage:              "Network-attached storage (NAS)",
		Action:             nasGatewayMain,
		CustomHelpTemplate: nasGatewayTemplate,
		HideHelp:           true,
	})
}

// Handler for 'obstor gateway nas' command line.
func nasGatewayMain(ctx *cli.Context) {
	// Validate gateway arguments.
	if !ctx.Args().Present() || ctx.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(ctx, obstor.NASBackendGateway, 1)
	}

	obstor.StartGateway(ctx, &NAS{ctx.Args().First()})
}

// NAS implements Gateway.
type NAS struct {
	path string
}

// Name implements Gateway interface.
func (g *NAS) Name() string {
	return obstor.NASBackendGateway
}

// NewGatewayLayer returns nas gatewaylayer.
func (g *NAS) NewGatewayLayer(creds auth.Credentials) (obstor.ObjectLayer, error) {
	var err error
	newObject, err := obstor.NewFSObjectLayer(g.path)
	if err != nil {
		return nil, err
	}
	return &nasObjects{newObject}, nil
}

// Production - nas gateway is production ready.
func (g *NAS) Production() bool {
	return true
}

// IsListenSupported returns whether listen bucket notification is applicable for this gateway.
func (n *nasObjects) IsListenSupported() bool {
	return false
}

func (n *nasObjects) StorageInfo(ctx context.Context) (si obstor.StorageInfo, _ []error) {
	si, errs := n.ObjectLayer.StorageInfo(ctx)
	si.Backend.GatewayOnline = si.Backend.Type == madmin.FS
	si.Backend.Type = madmin.Gateway
	return si, errs
}

// nasObjects implements gateway for Obstor and S3 compatible object storage servers.
type nasObjects struct {
	obstor.ObjectLayer
}

func (n *nasObjects) IsTaggingSupported() bool {
	return true
}
