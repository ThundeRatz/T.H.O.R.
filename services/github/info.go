package github

import (
	"thunderatz.org/thor/pkg/gclient"
)

// GetRepositories returns a slice of all repository names for the organization
func (ghs *Service) GetRepositories() []string {
	client, err := gclient.NewInstallationClient(ghs.AppID, ghs.InstallationID, ghs.PEMFile, &ghs.logger)

	if err != nil {
		ghs.logger.Error().Err(err).Msg("NewInstallationClient")
		return []string{}
	}

	return client.GetRepositories("ThundeRatz")
}

// GetMembers returns a slice of all members users for the organization
func (ghs *Service) GetMembers() []string {
	client, err := gclient.NewInstallationClient(ghs.AppID, ghs.InstallationID, ghs.PEMFile, &ghs.logger)

	if err != nil {
		ghs.logger.Error().Err(err).Msg("NewInstallationClient")
		return []string{}
	}

	return client.GetMembers("ThundeRatz")
}

// GetStats returns contributor statistics for all repositories in the organization
func (ghs *Service) GetStats() gclient.RepoStats {
	client, err := gclient.NewInstallationClient(ghs.AppID, ghs.InstallationID, ghs.PEMFile, &ghs.logger)

	if err != nil {
		ghs.logger.Error().Err(err).Msg("NewInstallationClient")
		return gclient.RepoStats{}
	}

	return client.GetStats("ThundeRatz")
}
