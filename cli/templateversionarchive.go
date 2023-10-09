package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/coder/pretty"
	"golang.org/x/xerrors"

	"github.com/coder/coder/v2/cli/clibase"
	"github.com/coder/coder/v2/cli/cliui"
	"github.com/coder/coder/v2/codersdk"
)

func (r *RootCmd) archiveTemplateVersion() *clibase.Cmd {
	var unarchive clibase.Bool

	client := new(codersdk.Client)
	cmd := &clibase.Cmd{
		Use:   "archive <template-name> [template-version-names...] ",
		Short: "Archive or unarchive a template version(s).",
		Middleware: clibase.Chain(
			r.InitClient(client),
		),
		Options: clibase.OptionSet{
			cliui.SkipPromptOption(),
			clibase.Option{
				Name:        "unarchive",
				Description: "Unarchive the selected template version.",
				Flag:        "unarchive",
				Value:       &unarchive,
			},
		},
		Handler: func(inv *clibase.Invocation) error {
			var (
				ctx      = inv.Context()
				versions []codersdk.TemplateVersion
			)

			organization, err := CurrentOrganization(inv, client)
			if err != nil {
				return err
			}

			if len(inv.Args) == 0 {
				return xerrors.Errorf("missing template name")
			}
			if len(inv.Args) < 2 {
				return xerrors.Errorf("missing template version name(s)")
			}

			templateName := inv.Args[0]
			template, err := client.TemplateByName(ctx, organization.ID, templateName)
			if err != nil {
				return xerrors.Errorf("get template by name: %w", err)
			}
			for _, versionName := range inv.Args[1:] {
				version, err := client.TemplateVersionByOrganizationAndName(ctx, organization.ID, template.Name, versionName)
				if err != nil {
					return xerrors.Errorf("get template version by name %q: %w", versionName, err)
				}
				versions = append(versions, version)
			}

			verb := "archived"
			if unarchive {
				verb = "unarchived"
			}
			failed := 0
			for _, version := range versions {
				if version.Archived == !unarchive.Value() {
					_, _ = fmt.Fprintln(
						inv.Stdout, fmt.Sprintf("Version "+pretty.Sprint(cliui.DefaultStyles.Keyword, version.Name)+" already "+verb),
					)
					continue
				}

				err := client.SetArchiveTemplateVersion(ctx, version.ID, !unarchive.Value())
				if err != nil {
					return xerrors.Errorf("set template version %q: %w", version.Name, err)
				}

				_, _ = fmt.Fprintln(
					inv.Stdout, fmt.Sprintf("Version "+pretty.Sprint(cliui.DefaultStyles.Keyword, version.Name)+" "+verb+" at "+cliui.Timestamp(time.Now())),
				)
			}

			if failed > 0 {
				return xerrors.Errorf("failed on %d template versions", failed)
			}
			return nil
		},
	}

	return cmd
}

func (r *RootCmd) archiveTemplateVersions() *clibase.Cmd {
	var all clibase.Bool
	client := new(codersdk.Client)
	cmd := &clibase.Cmd{
		Use:   "archive [template-name...] ",
		Short: "Archive unused failed template versions from a given template(s)",
		Middleware: clibase.Chain(
			r.InitClient(client),
		),
		Options: clibase.OptionSet{
			cliui.SkipPromptOption(),
			clibase.Option{
				Name:        "all",
				Description: "Include all unused template versions. By default, only failed template versions are archived.",
				Flag:        "all",
				Value:       &all,
			},
		},
		Handler: func(inv *clibase.Invocation) error {
			var (
				ctx           = inv.Context()
				templateNames = []string{}
				templates     = []codersdk.Template{}
			)

			organization, err := CurrentOrganization(inv, client)
			if err != nil {
				return err
			}

			if len(inv.Args) > 0 {
				templateNames = inv.Args

				for _, templateName := range templateNames {
					template, err := client.TemplateByName(ctx, organization.ID, templateName)
					if err != nil {
						return xerrors.Errorf("get template by name: %w", err)
					}
					templates = append(templates, template)
				}
			} else {
				template, err := selectTemplate(inv, client, organization)
				if err != nil {
					return err
				}

				templates = append(templates, template)
				templateNames = append(templateNames, template.Name)
			}

			// Confirm archive of the template.
			_, err = cliui.Prompt(inv, cliui.PromptOptions{
				Text:      fmt.Sprintf("Archive template versions of these templates: %s?", pretty.Sprint(cliui.DefaultStyles.Code, strings.Join(templateNames, ", "))),
				IsConfirm: true,
				Default:   cliui.ConfirmNo,
			})
			if err != nil {
				return err
			}

			failed := 0
			for _, template := range templates {
				resp, err := client.ArchiveTemplateVersions(ctx, template.ID, all.Value())
				if err != nil {
					return xerrors.Errorf("archive template %q: %w", template.Name, err)
				}

				_, _ = fmt.Fprintln(
					inv.Stdout, fmt.Sprintf("Archived %d versions from "+pretty.Sprint(cliui.DefaultStyles.Keyword, template.Name)+" at "+cliui.Timestamp(time.Now()), len(resp.ArchivedIDs)),
				)

				if ok, _ := inv.ParsedFlags().GetBool("verbose"); err == nil && ok {
					data, err := json.Marshal(resp)
					if err != nil {
						return xerrors.Errorf("marshal verbose response: %w", err)
					}
					_, _ = fmt.Fprintln(
						inv.Stdout, string(data),
					)
				}
			}

			if failed > 0 {
				return xerrors.Errorf("failed on %d templates", failed)
			}
			return nil
		},
	}

	return cmd
}
