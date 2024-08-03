#!groovy

def art = Artifactory.newServer(url: "https://art.maratos.tech", credentialsId: 'boo-sa-dba')

build = [
    app: "vip_patroni",
    version: "",
]

node('dockerhost') {
    timestamps {
        ansiColor('xterm') {
            try {
                stage('prepare') {

                    checkout scm

                    echo "Tag: ${env.TAG_NAME}"

                    if (env.TAG_NAME) {
                        build.version = env.TAG_NAME
                        if (build.version ==~ "v.*") {
                            def finder = (build.version =~ /v(.*)/)
                            build.version = finder.getAt(0)[1]
                        }
                    }

                    echo "Build env: ${build}"
                }
    
                stage('build') {
                    if (env.TAG_NAME) {
    
                        def build_image = docker.build("dba-vip_patroni", "-f Dockerfile . --build-arg app_version=${build.version}")
                        build_image.inside ("-v ${WORKSPACE}/build:/root/mnt"){
                            sh """
                            ls -al /root/build
                            cp /root/build/vip_patroni* /root/mnt/${build.version}.vip_patroni
                            """
                        }
                    }
                }
    
                stage('test') {
                    if (env.TAG_NAME) {
                        echo "test skip"
                    }
                }
    
               stage('upload') {
                    if (env.TAG_NAME) {
                        def uploadSpec = """{
                           "files": [
                               {
                                  "pattern": "${WORKSPACE}/build/*vip_patroni*",
                                  "target": "generic-local-dba/vip_patroni/"
                                }]
                            }"""
                        art.upload spec: uploadSpec, failNoOp: true
                    }
                }
            } catch(e) {
                currentBuild.result = 'FAILURE'
                throw(e)
            } finally {
                stage('cleanup') {
                    sh "docker rmi --force dba-vip_patroni"
        
                    deleteDir()
                }
            }
        }
    }
}