import os
import json

def dir_size(dir):
    for policy_type in os.listdir(dir):
        with open(f"docs/policies/{policy_type}.md", "w") as f:
            f.write(f"\n")
            for resource_type in os.listdir(os.path.join(dir,policy_type)):
                f.write(f"### {resource_type}\n")
                f.write("| Category | Resource | Severity | Description | Reference ID |\n")
                f.write("| -------- | -------- | -------- | ----------- | ------------ |\n")
                for (dirpath, dirs, files) in os.walk(os.path.join(dir, policy_type, resource_type)):
                    for filename in files:
                        if '.json' in filename:
                            with open(os.path.join(dirpath,filename)) as p:
                                policy = json.load(p)
                                category = policy['category']
                                resource = filename.split('.')[1]
                                severity = policy['severity']
                                description = policy['description'].replace('\n','')
                                reference_id = policy['reference_id']
                                f.write(f"| {category} | {resource} | {severity} | {description} | {reference_id} |\n")
                f.write("\n\n")

if __name__ == '__main__':
    policy_dir = os.path.join(os.getcwd(), "pkg", "policies", "opa", "rego")
    dir_size(policy_dir)
