# pithyblog
Lightweight blog

### Install and depoly

1. go get -v github.com/sryanyuan/pithyblog
2. go install github.com/sryanyuan/pithyblog
3. cd bin
4. Copy the config.example.toml to bin, and modify the config file as you want
5. pithyblog setup -config config.example.toml, and you will get the admin account and password, you must remember it.After you login the site,you can modify it.
6. Copy template and static directories to pithyblog directory
7. pithyblog run -config config.example.toml, now you can visit your site in your browser.