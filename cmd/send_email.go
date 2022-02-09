package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/curiosinauts/platformctl/pkg/mailutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed templates/*
var templates embed.FS

// emailCmd represents the email command
var emailCmd = &cobra.Command{
	Use:     "email",
	Short:   "sends email",
	Long:    `sends email`,
	PreRunE: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fromName := "Captain Curiosity"
		fromEmail := "captain@curiosityworks.org"
		subject := "Welcome"
		toName := ""
		toEmail := args[0]
		custom := "Learn about variables"
		plainTextContent := ""
		apiKey := viper.Get("sendgrid_api_key").(string)

		eh := ErrorHandler{"sending email"}

		t, err := template.New("pages").ParseFS(templates, "templates/*.html")
		eh.HandleError("getting email templates", err)

		buffer := bytes.NewBufferString("")
		err = t.ExecuteTemplate(buffer, "welcome.html", custom)
		eh.HandleError("rendering template", err)

		if debug {
			fmt.Println("email content:", buffer.String())
		}

		htmlContent := buffer.String()

		if false {
			response, emailErr := mailutil.Send(fromName, fromEmail, subject, toName, toEmail, plainTextContent, htmlContent, apiKey)
			eh.HandleError("email trasmission", emailErr)

			if debug {
				fmt.Println("header      :", response.Headers)
				fmt.Println("status code :", response.StatusCode)
				fmt.Println("body        :", response.Body)
			}
		}

		msg.Success("sending welcome email")
	},
}

func init() {
	sendCmd.AddCommand(emailCmd)
}
