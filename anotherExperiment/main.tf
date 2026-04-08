

module "iam" {
    source = "./modules/iam"

}

module "sqs"{
    source = "./modules/sqs"
}

module "lambda" {
    source = "./modules/lambda"
}