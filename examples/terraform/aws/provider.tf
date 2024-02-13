provider "aws" {
    region                      = "us-east-1"

    # The below allow the example to work with localstack.
    # You would remove this for a real implementation.
    s3_force_path_style         = true

    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true

    endpoints {
        s3 = "http://localhost:4566"
    }
}