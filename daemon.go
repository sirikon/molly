package main

import (
	"fmt"
	"mime/multipart"
	"github.com/urfave/cli"
	"github.com/gin-gonic/gin"
)

func ActionDaemon(c *cli.Context) error {
	router := gin.Default()

	router.POST("/deploy", func(c *gin.Context) {

		project := &Project{}

		if err := GetProjectByName(c.PostForm("project"), project); err != nil {
			c.String(400, "Specified project doesn't exist\n")
			fmt.Println(err)
			return
		}

		token := c.PostForm("token")

		if !project.CheckToken(token) {
			c.String(400, "Wrong token\n")
			return
		}

		var uploadedFile multipart.File

		if file, _, err := c.Request.FormFile("artifact"); err != nil {
			c.String(500, "Error reading the uploaded file\n")
			fmt.Println(err)
			return
		} else {
			uploadedFile = file
		}

		if err := project.StoreArtifact(uploadedFile); err != nil {
			c.String(500, "Error storing the artifact\n")
			fmt.Println(err)
			return
		}

		if err := project.RunDeploymentScript(); err != nil {
			c.String(500, "Error running project deploy script\n")
			fmt.Println(err)
			return
		}

		if err := project.RestartService(); err != nil {
			c.String(500, "Error restarting service\n")
			fmt.Println(err)
			return
		}

		c.String(200, "Done\n")
	})

	router.Run(":8080")

	return nil
}

