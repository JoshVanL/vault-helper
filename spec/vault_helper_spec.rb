require 'openssl'
require "docker"
require 'logger'

$logger = Logger.new(STDERR)
$logger.level = Logger::DEBUG
$containers = []

class VaultServer
  def image
      @image ||= Docker::Image.build_from_dir('.')
  end

  def start
    if @container.nil?
      @container = image.run('dev-server')
      $logger.info "waiting for server to be fully provisoned"
    end
    while true
      break if provisoned?
      sleep 1
    end
  end

  def get(path)
    uri = URI.parse(File.join(url,path))
    http = Net::HTTP.new(uri.host, uri.port)

    req =Net::HTTP::Get.new(uri.request_uri)
    req.add_field("X-Vault-Token", "root-token")

    http.request(req)
  end

  def provisoned?
    response = get('v1/auth/token/roles/cluster1-etcd')
    response.code == "200"
  rescue
    false
  end

  def url
    "http://#{@container.json['NetworkSettings']['IPAddress']}:8200"
  end

  def cleanup
    $logger.info "cleanup"
    unless @container.nil?
      @container.kill
      @container.wait
      @container.remove
    end
  end
end

$server = VaultServer.new
at_exit do
  $server.cleanup
  $containers.each do |container|
    container.kill
    container.wait
    container.remove
  end
end

describe "docker image" do
  let (:image) do
    $server.image
  end

  let (:vault_addr) do
    $server.start
    $server.url
  end

  let (:container) do
      container = Docker::Container.create(
        'Image' => image.id,
        'Cmd' => cmd,
        'Env' => environment,
      )
      $containers << container
      container.store_file("/etc/vault/environment", "VAULT_ADDR=#{vault_addr}")
      container.store_file("/etc/vault/token", "root-token")
      container.start
      container
  end

  describe "cert" do
    let :cmd do
      [
          'cert',
          '/tmp/test',
      ]
    end

    let :environment do
      [
          'VAULT_CERT_CN=kube-apiserver',
          'VAULT_CERT_ROLE=cluster1/pki/k8s/sign/kube-apiserver',
      ]
    end

    it "retrieves certificate with SANs" do
      expect(container.wait['StatusCode']).to eq(0), "expected successful execute: error stdout=#{container.logs(stdout: true)} stderr=#{container.logs(stderr: true)}"
      cert = OpenSSL::X509::Certificate.new container.read_file('/tmp/test.pem')
      key = OpenSSL::PKey::RSA.new container.read_file('/tmp/test-key.pem')
      expect(cert.check_private_key(key)).to eq(true), "Certificate is not matching key"

      ca = OpenSSL::X509::Certificate.new container.read_file('/tmp/test-ca.pem')
      store = OpenSSL::X509::Store.new
      store.add_cert ca
      expect(store.verify(cert)).to eq(true), "Certificate is not being verified by CA"
    end
  end
end

