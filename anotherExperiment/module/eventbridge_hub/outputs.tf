

# output "file_created_rule_arn" {
#     value = module.eventbridge.eventbridge_rule_arns["csv-created"]
# }

output "file_created_rule_arn" {
    value = module.eventbridge.eventbridge_rule_arns["batch-ready"]
}