package accurics

{{.prefix}}notEncryptedSageMaker[sgm_notebook.id] {
	sgm_notebook := input.aws_sagemaker_notebook_instance[_]
    object.get(sgm_notebook.config, "kms_key_id", "undefined") == [null, "undefined"][_]
}