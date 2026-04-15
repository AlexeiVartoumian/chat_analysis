resource "aws_lambda_layer_version" "requests_layer" {

    filename         = data.archive_file.requests_layer.output_path
    source_code_hash = data.archive_file.requests_layer.output_base64sha256
    layer_name = "requests_layer_data_source"
    compatible_runtimes = ["python3.14", "python3.13", "python3.12", "python3.11", "python3.10"]
}