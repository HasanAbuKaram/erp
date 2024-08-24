# erp

To run the docker:
$ docker-compose up --build


Yes, creating a new branch for developing and integrating the static file server is a good practice. This approach helps you manage changes more effectively and keeps your main branch stable. Here's how you can do it:

### 1. **Create a New Branch**

```bash
# Make sure you're in your project root directory
git checkout -b static-server
```

### 2. **Develop and Test the Static File Server**

Follow the steps outlined earlier to set up the static file server, Dockerfile, and `docker-compose.yml` changes. Make sure everything works as expected.

### 3. **Commit Changes**

Once you have completed the development and testing, commit your changes to the new branch:

```bash
git add .
git commit -m "Add static file server service"
```

### 4. **Push the Branch**

Push your new branch to the remote repository:

```bash
git push -u origin static-server
```

### 5. **Create a Pull Request (PR)**

- Go to your repository on GitHub.
- You should see an option to create a Pull Request for your new branch.
- Create a PR to merge `static-server` into `main`.
- Review the changes, and if everything looks good, merge the PR.

### 6. **Merge and Clean Up**

Once the PR is merged:

- Check out the `main` branch:

  ```bash
  git checkout main
  ```

- Pull the latest changes:

  ```bash
  git pull origin main
  ```

- Optionally, delete the `static-server` branch locally and remotely if it’s no longer needed:

  ```bash
  git branch -d static-server
  git push origin --delete static-server
  ```

Delete the master branch both locally and remotely:

Locally:
$ git branch -d master

Remotely:
$ git push origin --delete master

Creating a new branch and following these steps will help you maintain a clean workflow and ensure that new features are well-integrated into your main project.