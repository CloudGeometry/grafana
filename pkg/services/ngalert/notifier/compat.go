package notifier

import (
	"encoding/json"

	"github.com/prometheus/alertmanager/config"

	alertingNotify "github.com/grafana/alerting/notify"
	alertingTemplates "github.com/grafana/alerting/templates"

	"github.com/grafana/grafana/pkg/components/simplejson"
	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
)

func PostableGrafanaReceiverToGrafanaIntegrationConfig(p *apimodels.PostableGrafanaReceiver) *alertingNotify.GrafanaIntegrationConfig {
	return &alertingNotify.GrafanaIntegrationConfig{
		UID:                   p.UID,
		Name:                  p.Name,
		Type:                  p.Type,
		DisableResolveMessage: p.DisableResolveMessage,
		Settings:              json.RawMessage(p.Settings),
		SecureSettings:        p.SecureSettings,
	}
}

func PostableApiReceiverToApiReceiver(r *apimodels.PostableApiReceiver) *alertingNotify.APIReceiver {
	integrations := alertingNotify.GrafanaIntegrations{
		Integrations: make([]*alertingNotify.GrafanaIntegrationConfig, 0, len(r.GrafanaManagedReceivers)),
	}
	for _, cfg := range r.GrafanaManagedReceivers {
		integrations.Integrations = append(integrations.Integrations, PostableGrafanaReceiverToGrafanaIntegrationConfig(cfg))
	}

	return &alertingNotify.APIReceiver{
		ConfigReceiver:      r.Receiver,
		GrafanaIntegrations: integrations,
	}
}

func PostableApiAlertingConfigToApiReceivers(c apimodels.PostableApiAlertingConfig) []*alertingNotify.APIReceiver {
	apiReceivers := make([]*alertingNotify.APIReceiver, 0, len(c.Receivers))
	for _, receiver := range c.Receivers {
		apiReceivers = append(apiReceivers, PostableApiReceiverToApiReceiver(receiver))
	}
	return apiReceivers
}

type DecryptFn = func(value string) string

func PostableToGettableGrafanaReceiver(r *apimodels.PostableGrafanaReceiver, provenance *models.Provenance, decryptFn DecryptFn) (apimodels.GettableGrafanaReceiver, error) {
	out := apimodels.GettableGrafanaReceiver{
		UID:                   r.UID,
		Name:                  r.Name,
		Type:                  r.Type,
		DisableResolveMessage: r.DisableResolveMessage,
		SecureFields:          make(map[string]bool, len(r.SecureSettings)),
	}
	if provenance != nil {
		out.Provenance = apimodels.Provenance(*provenance)
	}

	if r.Settings == nil && r.SecureSettings == nil {
		return out, nil
	}

	settings := simplejson.New()
	if r.Settings != nil {
		var err error
		settings, err = simplejson.NewJson(r.Settings)
		if err != nil {
			return apimodels.GettableGrafanaReceiver{}, err
		}
	}

	for k, v := range r.SecureSettings {
		decryptedValue := decryptFn(v)
		if decryptedValue == "" {
			continue
		} else {
			settings.Set(k, decryptedValue)
		}
		out.SecureFields[k] = true
	}

	jsonBytes, err := settings.MarshalJSON()
	if err != nil {
		return apimodels.GettableGrafanaReceiver{}, err
	}

	out.Settings = jsonBytes

	return out, nil
}

func PostableToGettableApiReceiver(r *apimodels.PostableApiReceiver, provenances map[string]models.Provenance, decryptFn DecryptFn) (apimodels.GettableApiReceiver, error) {
	out := apimodels.GettableApiReceiver{
		Receiver: config.Receiver{
			Name: r.Receiver.Name,
		},
	}

	for _, gr := range r.GrafanaManagedReceivers {
		var prov *models.Provenance
		if p, ok := provenances[gr.UID]; ok {
			prov = &p
		}

		gettable, err := PostableToGettableGrafanaReceiver(gr, prov, decryptFn)
		if err != nil {
			return apimodels.GettableApiReceiver{}, err
		}
		out.GrafanaManagedReceivers = append(out.GrafanaManagedReceivers, &gettable)
	}

	return out, nil
}

// ToTemplateDefinitions converts the given PostableUserConfig's TemplateFiles to a slice of TemplateDefinitions.
func ToTemplateDefinitions(cfg *apimodels.PostableUserConfig) []alertingTemplates.TemplateDefinition {
	out := make([]alertingTemplates.TemplateDefinition, 0, len(cfg.TemplateFiles))
	for name, tmpl := range cfg.TemplateFiles {
		out = append(out, alertingTemplates.TemplateDefinition{
			Name:     name,
			Template: tmpl,
		})
	}
	return out
}

// Silence-specific compat functions to convert between grafana/alerting and model types.

func GettableSilenceToSilence(s alertingNotify.GettableSilence) *models.Silence {
	sil := models.Silence(s)
	return &sil
}

func GettableSilencesToSilences(silences alertingNotify.GettableSilences) []*models.Silence {
	res := make([]*models.Silence, 0, len(silences))
	for _, sil := range silences {
		res = append(res, GettableSilenceToSilence(*sil))
	}
	return res
}

func SilenceToPostableSilence(s models.Silence) *alertingNotify.PostableSilence {
	var id string
	if s.ID != nil {
		id = *s.ID
	}
	return &alertingNotify.PostableSilence{
		ID:      id,
		Silence: s.Silence,
	}
}
