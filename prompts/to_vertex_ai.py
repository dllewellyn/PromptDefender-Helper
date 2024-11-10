import yaml
import json
import re
import argparse


def parse_scoring_prompt(file_path):
    with open(file_path, 'r') as file:
        content = file.read()

    # Extract frontmatter parameters and model
    frontmatter_pattern = re.compile(
        r'^---\s*\n(.*?)\n---', re.DOTALL | re.MULTILINE)
    frontmatter_match = frontmatter_pattern.search(content)
    if frontmatter_match:
        frontmatter_yaml = frontmatter_match.group(1)
        frontmatter_data = yaml.safe_load(frontmatter_yaml)
        config = frontmatter_data.get('config', {})
        parameters = {k: config[k] for k in [
            'temperature', 'topK', 'topP', 'tokenLimits'] if k in config}
        model = frontmatter_data.get('model', '')
    else:
        parameters = {}
        model = ''

    # Remove frontmatter from content
    content = content[frontmatter_match.end(
    ):] if frontmatter_match else content

    # Extract system prompt
    system_role_pattern = re.compile(
        r'\{\{role\s+"system"\}\}(.*?)(?=(\{\{role|\Z))', re.DOTALL)
    system_role_match = system_role_pattern.search(content)
    if system_role_match:
        prompt_text = system_role_match.group(1).strip()
        content = content.replace(system_role_match.group(0), '')
    else:
        prompt_text = ''

    # Extract examples from user and model roles
    examples = []
    example_pattern = re.compile(
        r'\{\{role\s+"user"\}\}\s*(.*?)\s*\{\{role\s+"model"\}\}\s*(.*?)(?=(\{\{role|\Z))', re.DOTALL)
    matches = example_pattern.findall(content)
    for user_text, model_text, _ in matches:
        examples.append({
            "inputs": [user_text.strip()],
            "outputs": [model_text.strip()]
        })

    # Create final JSON structure
    output = {
        "parameters": parameters,
        "model": model,
        "type": "freeform",
        "systemInstruction": {
            "parts": [
                {
                    "text": prompt_text
                }
            ]
        },
        "examples": examples
    }

    return output


def main():
    parser = argparse.ArgumentParser(description='Convert scoring prompt to JSON format.')
    parser.add_argument('input_file', type=str, help='Path to the input scoring prompt file')
    parser.add_argument('output_file', type=str, help='Path to save the converted JSON output')

    args = parser.parse_args()

    data = parse_scoring_prompt(args.input_file)

    with open(args.output_file, 'w') as f:
        json.dump(data, f, indent=2)

    print(f"Conversion complete. Output saved to {args.output_file}")


if __name__ == "__main__":
    main()
