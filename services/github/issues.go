package github

import "thunderatz.org/thor/pkg/gclient"

// SendDefaultVSSIssueMessage send the default message to newly created issues
func (ghs *Service) SendDefaultVSSIssueMessage(issueNumber int) {
	client, err := gclient.NewInstallationClient(ghs.AppID, ghs.InstallationID, ghs.PEMFile, &ghs.logger)

	if err != nil {
		ghs.logger.Error().Err(err).Msg("NewInstallationClient")
		return
	}

	client.IssueComment("ThundeRatz", "vss_simulation", `Hi! Thank you for opening an issue for this project!  \
Please, make sure you followed the project's [contribution guidelines](https://github.com/ThundeRatz/vss_simulation/blob/feature/open_source/CONTRIBUTING.md), a team member will answer when possible!

--

Olá! Obrigado por abrir uma isse para esse projeto!  \
Por favor, tenha certeza que leu  as [diretrizes de contribuição](https://github.com/ThundeRatz/vss_simulation/blob/feature/open_source/CONTRIBUTING.pt-br.md) do projeto, alguém da equipe responderá assim que possível!`, issueNumber)
}
