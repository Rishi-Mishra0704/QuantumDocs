package template

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
)

var (
	HtmlContent string
	HtmlMu      sync.RWMutex
	tmpl        *template.Template
)

func init() {
	var err error
	tmpl, err = template.New("apidoc").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <title>{{.Title}} API Documentation</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
        }
        h1, h2, h3 {
            color: #333;
        }
        .endpoint {
            background-color: #f4f4f4;
            border: 1px solid #ddd;
            border-radius: 5px;
            margin-bottom: 20px;
            padding: 15px;
        }
        .method {
            font-weight: bold;
            color: #0066cc;
        }
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
        .schema {
            background-color: #f8f8f8;
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 15px;
            white-space: pre-wrap;
            font-family: monospace;
        }
    </style>
<script>
    console.log("WebSocket script loaded");
    var socket = new WebSocket("ws://localhost:8080/ws");

    socket.onmessage = function (event) {
        console.log("Message received from server:", event.data);
        if (event.data === "reload") {
            location.reload();
        }
    };

    socket.onclose = function (event) {
        console.log("WebSocket closed. Reconnecting...");
        setTimeout(function() {
            socket = new WebSocket("ws://localhost:8080/ws");
        }, 5000);
    };

    socket.onerror = function (error) {
        console.error("WebSocket error:", error);
    };
</script>

</head>
<body>
    <h1>{{.Title}} API Documentation</h1>
    <p>{{.Description}}</p>
    <p>Version: {{.Version}}</p>
    {{range .Endpoints}}
    <div class="endpoint">
        <h2><span class="method">{{.Method}}</span> {{.Path}}</h2>
        <p>{{.Description}}</p>
        {{if .Parameters}}
        <div>
            <h3>Parameters</h3>
            <table>
                <tr>
                    <th>Name</th>
                    <th>Type</th>
                    <th>Required</th>
                    <th>Description</th>
                </tr>
                {{range .Parameters}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.Type}}</td>
                    <td>{{.Required}}</td>
                    <td>{{.Description}}</td>
                </tr>
                {{end}}
            </table>
        </div>
        {{end}}
        {{if .RequestSchema}}
        <div>
            <h3>Request Schema</h3>
            <pre class="schema">{{.RequestSchema}}</pre>
        </div>
        {{end}}
        {{if .ResponseSchema}}
        <div>
            <h3>Response Schema</h3>
            <pre class="schema">{{.ResponseSchema}}</pre>
        </div>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`)
	if err != nil {
		panic(err)
	}
}

// GenerateHTML generates the HTML content for the API documentation
func GenerateHTML(doc *models.APIDoc) string {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, doc)
	if err != nil {
		return "Error generating HTML: " + err.Error()
	}
	return buf.String()
}
