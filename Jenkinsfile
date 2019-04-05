def label = "buildpod.${env.JOB_NAME}".replaceAll(/[^A-Za-z-]+/, '-').take(62) + "p"

podTemplate(
	name: label,
	label: label,
	containers: [
		containerTemplate(
			name: 'build-golang',
			image: 'eu.gcr.io/smedia-kubernetes/build-golang:1.12.0-1.3.2',
			ttyEnabled: true,
			command: 'cat',
			resourceRequestCpu: '500m',
			resourceRequestMemory: '750Mi',
			resourceLimitCpu: '2000m',
			resourceLimitMemory: '750Mi',
		),
	],
	volumes: [],
	inheritFrom: '',
	namespace: 'jenkins',
	serviceAccount: '',
	workspaceVolume: emptyDirWorkspaceVolume(false),
) {
	node(label) {
		properties([
			buildDiscarder(logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '3', numToKeepStr: '5')),
			pipelineTriggers([
				cron('H 2 * * *'),
				pollSCM('H/5 * * * *'),
			]),
		])
		try {
			container('build-golang') {
				stage('Golang Checkout') {
					timeout(time: 5, unit: 'MINUTES') {
						checkout([
							$class: 'GitSCM',
							branches: scm.branches,
							doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
							extensions: scm.extensions + [[$class: 'CloneOption', noTags: false, reference: '', shallow: false]],
							submoduleCfg: [],
							userRemoteConfigs: scm.userRemoteConfigs
						])
					}
				}
				stage('Golang Link') {
					timeout(time: 5, unit: 'MINUTES') {
						sh """
						mkdir -p /go/src/github.com/seibert-media
						ln -s `pwd` /go/src/github.com/seibert-media/aptly-utils
						"""
					}
				}
				stage('Golang Test') {
					timeout(time: 15, unit: 'MINUTES') {
						sh "cd /go/src/github.com/seibert-media/aptly-utils && make test"
					}
				}
			}
			currentBuild.result = 'SUCCESS'
		} catch (any) {
			currentBuild.result = 'FAILURE'
			throw any //rethrow exception to prevent the build from proceeding
		} finally {
			if ('FAILURE'.equals(currentBuild.result)) {
				emailext(
					body: '${DEFAULT_CONTENT}',
					mimeType: 'text/plain',
					replyTo: '$DEFAULT_REPLYTO',
					subject: '${DEFAULT_SUBJECT}',
					recipientProviders:[
						[$class: 'CulpritsRecipientProvider'],
						[$class: 'RequesterRecipientProvider']
					]
				)
			}
		}
	}
}