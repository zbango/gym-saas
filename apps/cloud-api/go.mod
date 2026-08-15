module github.com/zbango/gym-saas/apps/cloud-api

go 1.23.0

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/zbango/gym-saas/go/core v0.0.0
)

require github.com/stretchr/testify v1.10.0 // indirect

replace github.com/zbango/gym-saas/go/core => ../../go/core
