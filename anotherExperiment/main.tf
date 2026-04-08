

module "iam_hub" {
    source = "./modules/iam_hub"
    spoke_accounts = var.spoke_accounts
    hub_account = var.hub_account
}

module "sqs"{
    source = "./modules/sqs"
}

module "lambda" {
    source = "./modules/lambda"
}