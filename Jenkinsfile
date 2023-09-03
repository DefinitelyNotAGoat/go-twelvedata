pipeline {
    agent {
        kubernetes {
            defaultContainer 'jnlp'
            yamlFile 'agentpod.yml'
        }
    }
    stages {
        stage('Lint') {
            steps {
                container('golint') {
                    sh 'golangci-lint run --timeout 5m ./...'
                }
            }
        }
        stage('Test') {
            steps {
                container('golint') {
                    sh 'go test -v --run=Unit ./...'
                }
            }
        }
        stage('Integration Test') {
            when {
                expression { env.BRANCH_NAME == 'main' }
            }
            steps {
                container('golint') {
                    sh 'go test -v --run=Integration ./...'
                }
            }
        }
    }
}