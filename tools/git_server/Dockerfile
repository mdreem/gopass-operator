FROM ubuntu:20.04

RUN DEBIAN_FRONTEND="noninteractive" apt update && apt install  openssh-server sudo git -y --no-install-recommends
RUN groupadd git
RUN useradd -rm -d /home/git -s /bin/bash -g git -G sudo -u 1000 git
RUN echo 'git:git' | chpasswd
RUN service ssh start

EXPOSE 22

COPY password-store.tgz /home/git
RUN mkdir /home/git/password-store && cd /home/git/password-store && tar xzf ../password-store.tgz && rm ../password-store.tgz
RUN chmod +rw -R /home/git/

CMD ["/usr/sbin/sshd","-D"]
