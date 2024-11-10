import json
import yaml
import argparse

def convert_json_to_dotprompt(json_file, output_file):
    with open(json_file, 'r') as f:
        data = json.load(f)

    # Prepare frontmatter
    frontmatter = {
        'model': data.get('model', ''),
        'config': data.get('parameters', {}),
        'tools': []
    }

    # Start building the .prompt content
    dotprompt_content = []

    # Add frontmatter
    frontmatter_yaml = yaml.safe_dump(frontmatter, default_flow_style=False).strip()
    dotprompt_content.append('---')
    dotprompt_content.append(frontmatter_yaml)
    dotprompt_content.append('---\n')

    # Add system role
    system_instruction = ''
    system_parts = data.get('systemInstruction', {}).get('parts', [])
    for part in system_parts:
        system_instruction += part.get('text', '') + '\n'

    dotprompt_content.append('{{role "system"}}')
    dotprompt_content.append(system_instruction.strip() + '\n')

    # Add examples
    examples = data.get('examples', [])
    for example in examples:
        # User input
        user_input = example.get('inputs', [''])[0]
        dotprompt_content.append('{{role "user"}}')
        dotprompt_content.append(user_input.strip() + '\n')

        # Model output
        model_output = example.get('outputs', [''])[0]
        dotprompt_content.append('{{role "model"}}')
        dotprompt_content.append(model_output.strip() + '\n')

    # Write to output file
    with open(output_file, 'w') as f:
        for line in dotprompt_content:
            f.write(line + '\n')

    print(f"Conversion complete. Output saved to {output_file}")

def main():
    parser = argparse.ArgumentParser(description='Convert prompt_scoring.json to .prompt format.')
    parser.add_argument('input_file', type=str, help='Path to the input JSON file')
    parser.add_argument('output_file', type=str, help='Path to save the converted .prompt file')

    args = parser.parse_args()

    convert_json_to_dotprompt(args.input_file, args.output_file)

if __name__ == "__main__":
    main()