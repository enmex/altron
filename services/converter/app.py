from flask import Flask, request, jsonify
from generator import generate_sploit
import base64
import subprocess
import os
import shutil

app = Flask(__name__)

@app.route('/convert_to_curl', methods=['POST'])
def convert_to_curl():
    data = request.get_json()

    if 'base64_str' not in data:
        return jsonify({"error": "No base64_str field in the request"}), 400

    try:
        decoded_bytes = base64.b64decode(data['base64_str'])
        with open('output_file', 'wb') as f:
            f.write(decoded_bytes)

        cmd = '''./h2c < output_file'''
        output = subprocess.Popen(cmd,
                                  stdout=subprocess.PIPE,
                                  stderr=subprocess.STDOUT,
                                  shell=True,
                                  universal_newlines=True)
        message = output.stdout.read()
        output.terminate()
        return jsonify({"result": message})
    except Exception as e:
        return jsonify({"error": str(e)}), 400
    

@app.route('/convert_to_requests', methods=['POST'])
def convert_to_requests():
    data = request.get_json()

    if 'base64_str' not in data:
        return jsonify({"error": "No base64_str field in the request"}), 400

    try:
        decoded_bytes = base64.b64decode(data['base64_str'])
        with open('output_file', 'wb') as f:
            f.write(decoded_bytes)

        cmd = '''./h2c < output_file'''
        output = subprocess.Popen(cmd,
                                  stdout=subprocess.PIPE,
                                  stderr=subprocess.STDOUT,
                                  shell=True,
                                  universal_newlines=True)
        message = output.stdout.read()
        output.terminate()

        cmd = 'echo \'{}\' | curlconverter -'.format(message)
        output = subprocess.Popen(cmd,
                                  stdout=subprocess.PIPE,
                                  stderr=subprocess.STDOUT,
                                  shell=True,
                                  universal_newlines=True)
        message = output.stdout.read()
        output.terminate()

        return jsonify({"result": message})
    except Exception as e:
        return jsonify({"error": str(e)}), 400


@app.route('/convert_to_requests_sploit', methods=['POST'])
def convert_to_requests_sploit():
    data = request.get_json()

    if 'base64_strings' not in data:
        return jsonify({"error": "No base64_strings field in the request"}), 400

    try:
        packets = data['base64_strings']
        messages = []

        for packet in packets:

            decoded_bytes = base64.b64decode(packet)

            with open('output_file', 'wb') as f:
                f.write(decoded_bytes)

            cmd = '''./h2c < output_file'''
            output = subprocess.Popen(cmd,
                                      stdout=subprocess.PIPE,
                                      stderr=subprocess.STDOUT,
                                      shell=True,
                                      universal_newlines=True)
            message = output.stdout.read()
            output.terminate()

            cmd = 'echo \'{}\' | curlconverter -'.format(message)
            output = subprocess.Popen(cmd,
                                      stdout=subprocess.PIPE,
                                      stderr=subprocess.STDOUT,
                                      shell=True,
                                      universal_newlines=True)
            message = output.stdout.read()

            messages.append(message)
            output.terminate()

        sploit = generate_sploit(messages)

        return jsonify({"result": sploit})
    except Exception as e:
        return jsonify({"error": str(e)}), 400




if __name__ == '__main__':
    app.run(debug=True, host="0.0.0.0", port=64004)