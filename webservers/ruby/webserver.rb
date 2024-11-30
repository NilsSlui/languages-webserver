require 'webrick'
require 'digest'
require 'base64'
require 'cgi'

# Create a WEBrick server
server = WEBrick::HTTPServer.new(Port: 8008)

# Define the POST `/sha256` route
server.mount_proc '/sha256' do |req, res|
  if req.request_method == 'POST'
    input = req.query['input']
    if input.nil? || input.empty?
      res.status = 400
      res.body = 'Input is required'
    else
      res.body = Digest::SHA256.hexdigest(input)
    end
  else
    res.status = 405
    res.body = 'Method Not Allowed'
  end
end

# Define the POST `/base64` route
server.mount_proc '/base64' do |req, res|
  if req.request_method == 'POST'
    input = req.query['input']
    if input.nil? || input.empty?
      res.status = 400
      res.body = 'Input is required'
    else
      res.body = Base64.strict_encode64(input)
    end
  else
    res.status = 405
    res.body = 'Method Not Allowed'
  end
end

# Define the POST `/urlencode` route
server.mount_proc '/urlencode' do |req, res|
  if req.request_method == 'POST'
    input = req.query['input']
    if input.nil? || input.empty?
      res.status = 400
      res.body = 'Input is required'
    else
      res.body = CGI.escape(input)
    end
  else
    res.status = 405
    res.body = 'Method Not Allowed'
  end
end

# Define the GET `/` route
server.mount_proc '/' do |req, res|
  if req.request_method == 'GET'
    res.status = 201
    res.body = 'ruby'
  else
    res.status = 405
    res.body = 'Method Not Allowed'
  end
end

# Trap interrupt signal to gracefully stop the server
trap 'INT' do
  server.shutdown
end

# Start the server
server.start
