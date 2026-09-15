$wazuhEvent = '{"rule": {"id": "5710", "level": 12, "description": "SSH Brute Force", "groups": ["syslog", "sshd", "authentication_failed"]}, "agent": {"id": "001", "name": "webserver"}, "manager": {"name": "wazuh-manager"}, "data": {"srcip": "192.168.1.99", "dstuser": "root"}}'
$zeekEvent = '{"ts":1698246245.123456,"uid":"C5bQ8u1k8","id.orig_h":"192.168.1.99","id.orig_p":54321,"id.resp_h":"10.0.0.5","id.resp_p":22,"proto":"tcp","service":"ssh","duration":10.5,"orig_bytes":1500,"resp_bytes":2000,"conn_state":"SF"}'
$suricataEvent = '{"timestamp":"2023-10-25T15:04:05.123456Z","event_type":"alert","src_ip":"192.168.1.99","src_port":54321,"dest_ip":"10.0.0.5","dest_port":22,"proto":"TCP","alert":{"action":"allowed","gid":1,"signature_id":2010935,"rev":3,"signature":"ET SCAN Suspicious inbound to SSH port 22","category":"Attempted Information Leak","severity":2}}'

# We must append without locking the file exclusively, which out-file does by default in some versions, but cmd echo doesn't.
cmd /c "echo $wazuhEvent >> `"E:\New folder\SAE_Tools\SAE\engine\logs\alerts\alerts.json`""
cmd /c "echo $zeekEvent >> `"E:\New folder\SAE_Tools\SAE\engine\logs\zeek\conn.log`""
cmd /c "echo $suricataEvent >> `"E:\New folder\SAE_Tools\SAE\engine\logs\suricata\eve.json`""
