clean:
	./hack/virt-delete-aio.sh || true
	rm -rf mydir zezere
	docker rm -f zezere

generate:
	mkdir mydir
	cp ./install-config.yaml mydir/
	OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="quay.io/openshift-release-dev/ocp-release:4.7.0-fc.4-x86_64" ./bin/openshift-install create aio-config --dir=mydir

start:
	./hack/virt-install-aio-ign.sh ./mydir/aio.ign

network:
	./hack/virt-create-net.sh

serve-ign:
	mkdir -p zezere
	cp cp mydir/aio.ign ./zezere/52:54:00:ee:42:e1
	chmod 777 ./zezere/52:54:00:ee:42:e1
	docker run --name zezere -v `pwd`/zezere:/usr/share/nginx/html/zezere/netboot/x86_64/ignition:ro,Z -p 37507:80 -d nginx

ssh:
	ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@192.168.126.10

image:
	curl -O -L https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.6/46.82.202008181646-0/x86_64/rhcos-46.82.202008181646-0-qemu.x86_64.qcow2.gz
	mv rhcos-46.82.202008181646-0-qemu.x86_64.qcow2.gz /tmp
	sudo gunzip /tmp/rhcos-46.82.202008181646-0-qemu.x86_64.qcow2.gz
