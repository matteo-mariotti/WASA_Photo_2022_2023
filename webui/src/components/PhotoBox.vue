
<script>
import backend from "@/services/axios";

export default {
  props: ["imageData"],
  data: function () {
    return {
      loading: null,
      binaryImage: null,
      likeCount: null,
      commentCount: null,
      ready: null,
      liked: this.imageData.loggedLike,
      commentBox: null
      // datetime: this.imageData. TODO Aggiungere la data
    }
  },
  methods: {
    async loadImage() {
      try {
        this.ready = false
        let response = await backend.get(`/photos/${this.imageData.id}`, {responseType: 'arraybuffer'});

        console.log("Request executed");
        console.log(response)
        this.handleResponse(response);

        this.ready = true

      } catch (error) {
        console.log("Error getting image data")
        console.log(error)
        this.handleError(error);
      }
    },
    handleResponse(response) {
      this.errormsg = null

      // Create the element and the blob data for the image
      let img = document.getElementById(this.imageData.id)
      let blob = new Blob([response.data])
      img.src = URL.createObjectURL(blob)

      console.log("Blob created successfully")

    },
    handleError(error) {
      this.errormsg = error.response.data.message
      // Print the error from the server (for debugging)
      // Log the error on console
      console.log(error)
    },
    async addLike() {
      try {
        let response = await backend.put(`/users/${this.imageData.photoOwner}/photos/${this.imageData.id}/likes/${sessionStorage.getItem("token")}`); // TODO Sitemare i parametri
        this.liked = true
        this.imageData.likes = this.imageData.likes+1
      } catch (error) {
        this.handleError(error);
      }
    },
    async removeLike() {
      try {
        let response = await backend.delete (`/users/${this.imageData.photoOwner}/photos/${this.imageData.id}/likes/${sessionStorage.getItem("token")}`);
        this.liked = false
        this.imageData.likes = this.imageData.likes-1
      } catch (error) {
        this.handleError(error);
      }
    },
    async comments(photoID) {
      let response = await backend.get (`/photos/${photoID}/comments`);
      this.commentBox = true


    }
  },
  mounted() {
    this.loadImage();
  }
}
// TODO Questo file è ancora da fare!!
</script>


<template>

  <div class="row d-flex justify-content-center" style="background-color: greenyellow;padding: 30px">

    <img v-bind:id="this.imageData.id" class="rounded-5 img-fluid" alt="User image">

    <div class="row" style="background-color: #6f42c1; margin-top: 1vh">

      <div style="text-align: center;" class="col-sm">
        <div class="row" style="background-color: #20c997">
          <div class="col-sm-6 center" style="background-color: #0dcaf0; min-height: 7vh">
            <!-- Empty heart-->
            <svg v-if="!liked" @click="addLike()" xmlns="http://www.w3.org/2000/svg" width="100%" fill="currentColor"
                 class="bi bi-heart"
                 viewBox="0 0 16 16">
              <path
                  d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.01L8 2.748zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143c.06.055.119.112.176.171a3.12 3.12 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15z"/>
            </svg>
            <!-- Filled heart-->
            <svg v-if="liked" @click="removeLike()" style="fill: red" xmlns="http://www.w3.org/2000/svg" width="100%" class="bi bi-heart-fill" viewBox="0 0 16 16">
              <path fill-rule="evenodd" d="M8 1.314C12.438-3.248 23.534 4.735 8 15-7.534 4.736 3.562-3.248 8 1.314z"/>
            </svg>
          </div>
          <div class="col-sm-6 center" style="font-size: 6vh">
            {{this.imageData.likes}}
          </div>
        </div>

      </div>
      <div class="col-sm"></div>
      <div class="col-sm">
        <div class="row" style="background-color: #0c4128; text-align: center">
          <div class="col-sm-6 center" style="font-size: 6vh">
            {{this.imageData.comments}}
          </div>
          <div class="col-sm-6 center" style="background-color: #0dcaf0; min-width: 7vh">
            <!-- Comments icon-->
            <svg @click="comments(this.imageData.id)" xmlns="http://www.w3.org/2000/svg" width="100%" fill="currentColor" class="bi bi-chat-left-text"
                 viewBox="0 0 16 16">
              <path
                  d="M14 1a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H4.414A2 2 0 0 0 3 11.586l-2 2V2a1 1 0 0 1 1-1h12zM2 0a2 2 0 0 0-2 2v12.793a.5.5 0 0 0 .854.353l2.853-2.853A1 1 0 0 1 4.414 12H14a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H2z"/>
              <path
                  d="M3 3.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zM3 6a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9A.5.5 0 0 1 3 6zm0 2.5a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5z"/>
            </svg>
          </div>
        </div>
      </div>

    </div>
  </div>

</template>

<style>

.center {
  display: flex;
  justify-content: center;
  align-items: center;
  border: 3px solid green;
}

</style>
