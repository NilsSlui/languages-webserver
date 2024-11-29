from flask import Flask, request, make_response, abort
import hashlib
import base64
import urllib.parse

app = Flask(__name__)

@app.route('/', defaults={'path': ''}, methods=['GET', 'POST'])
@app.route('/<path:path>', methods=['GET', 'POST'])
def handle_request(path):
    try:
        method = request.method
        uri = '/' + path

        if method == 'POST':
            input_value = request.form.get('input')
            if input_value is None:
                raise Exception('Input is required')

            if uri == '/sha256':
                output = hashlib.sha256(input_value.encode()).hexdigest()
                return output
            elif uri == '/base64':
                output = base64.b64encode(input_value.encode()).decode()
                return output
            elif uri == '/urlencode':
                output = urllib.parse.quote(input_value)
                return output
            else:
                return 'Not Found', 404

        elif method == 'GET':
            response = make_response('python', 201)
            return response

        else:
            return 'Method Not Allowed', 405

    except Exception as e:
        return str(e) or 'Internal Server Error', 500

if __name__ == '__main__':
    app.run(port=8008)
