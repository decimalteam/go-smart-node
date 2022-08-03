package client

import (
	"bitbucket.org/decimalteam/go-smart-node/x/upgrade/client/cli"
	"bitbucket.org/decimalteam/go-smart-node/x/upgrade/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
var CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal, rest.ProposalCancelRESTHandler)
