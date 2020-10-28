node {

    def registry = '' 
    def registryCredential = ''
    def commit_id = ''
	def dockerImage = ''

    //agent { dockerfile true }
	stage('Clone repository') {
        /* Cloning the Repository to our Workspace */
        checkout scm
    }
	stage('Build image') {
        /* This builds the actual image; synonymous to
		* docker build on the command line */
        commit_id = sh(returnStdout: true, script: 'git rev-parse HEAD')
  		echo "$commit_id"
        dockerImage = docker.build registry

	}
	stage('Tag and Register image') {
	    /* Finally, we'll push the image with tags:
    	* First, the git commit id.
    	* Second, the app name with git commit.
    	* Third, latest tag.*/
        docker.withRegistry( '', registryCredential ) {
            dockerImage.push("$commit_id")
            dockerImage.push("webapp")
			dockerImage.push("webapp_${env.BUILD_NUMBER}")
            dockerImage.push("latest")
		}
    }
	stage('Remove Unused docker image') {
		/* Cleaning from local machine */
		//sh "docker rmi -f `docker images -q`"
		 sh "docker rmi -f `docker images | grep webapp | awk '{print \$3}'`"
	}
}
