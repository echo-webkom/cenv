# How to publish a new version

1. Tag your commit

   ```sh
   git tag <version>
   ```

1. Push tagged commit to origin

   ```sh
   git push
   ```

1. Create a new release in the repo with the title being the release version (eg. `v.1.5.2`). The workflow will start automatically when the release is created.
