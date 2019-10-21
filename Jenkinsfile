pipeline {
    agent any
    environment {
        BUILD_STATUS = 'success'
    }
    stages {
        stage {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'exit 1'
                }
                echo currentBuild.result
                echo currentStage.result
            }
        }
    }
}