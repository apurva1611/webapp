node {

    def dockerImage
    def registryCredential = 'DockerHub'
    def commit_id
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
        dockerImage = docker.build ("webapp", "-f Dockerfile .")

	}
	stage('Tag and Register image') {
	    /* Finally, we'll push the image with tags:
    	* First, the git commit id.
    	* Second, the app name with git commit.
    	* Third, latest tag.*/
        docker.withRegistry( '', registryCredential ) {
            dockerImage.push("$commit_id")
            dockerImage.push("cloud-webapp_$commit_id")
            dockerImage.push("latest")
		}
    }
	stage('Remove Unused docker image') {
		/* Cleaning from local machine */
		//sh "docker rmi -f `docker images -q`"
		 sh "docker rmi -f `docker images | grep webapp | awk '{print \$3}'`"
	}
}