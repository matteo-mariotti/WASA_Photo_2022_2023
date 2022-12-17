<script>
import backend from "@/services/axios";
import SuccessMsg from "@/components/SuccessMsg.vue";

export default {
  components: { SuccessMsg },
  props: ["imageData"],
  data: function () {
    return {
      ready: null,
      liked: this.imageData.loggedLike,
      isOwned: this.imageData.photoOwner === sessionStorage.getItem('token'),
      successmsg: null
    }
  },
  methods: {
    async loadImage() {
      try {
        this.ready = false
        let response = await backend.get(`/photos/${this.imageData.id}`, { responseType: 'arraybuffer' });

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
        await backend.put(`/users/${this.imageData.photoOwner}/photos/${this.imageData.id}/likes/${sessionStorage.getItem("token")}`); // TODO Sitemare i parametri
        this.liked = true
        this.imageData.likes = this.imageData.likes + 1
      } catch (error) {
        this.handleError(error);
      }
    },
    async removeLike() {
      try {
        await backend.delete(`/users/${this.imageData.photoOwner}/photos/${this.imageData.id}/likes/${sessionStorage.getItem("token")}`);
        this.liked = false
        this.imageData.likes = this.imageData.likes - 1
      } catch (error) {
        this.handleError(error);
      }
    },
    comments() {
      alert("Box dei commenti");
    },
    async deletePhoto() {
      try {
        await backend.delete(`/users/${this.imageData.photoOwner}/photos/${this.imageData.id}`);
        this.successmsg = "Photo deleted successfully";
        setTimeout(this.reloadPage, 1500);
      } catch (error) {

        this.handleError(error);
      }
    },
    reloadPage() {
      location.reload();
    }
  },
  mounted() {
    this.loadImage();
  }
}

</script>

<template>

  <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

  <div class="card m-3" style="width: 25rem; background-color: #FFFFF0">
    <div class="card-body">
      <div class="d-flex">
        <h5 class="card-title">{{ this.imageData.photoOwner }}</h5>
        <button @click="deletePhoto()" v-if="isOwned" class="btn btn-outline-danger btn-sm ms-auto"> Delete photo
        </button>
      </div>
      <p class="card-text" style="text-align: right">{{ this.imageData.date }}</p>
    </div>
    <img class="card-img ps-4 pe-4" v-bind:id=this.imageData.id alt="Card image cap">
    <div class="d-flex card-body justify-content-evenly">
      <span>
        <svg v-if="!liked" @click="addLike()" viewBox="0 0 24 24" width="24" height="24" stroke="currentColor"
          stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" class="css-i6dzq1">
          <path
            d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z">
          </path>
        </svg>
        <svg v-if="liked" @click="removeLike()" style="fill: red" viewBox="0 0 24 24" width="24" height="24"
          stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"
          class="css-i6dzq1">
          <path
            d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z">
          </path>
        </svg>

        {{ this.imageData.likes }}
      </span>
      <span>
        <svg @click="comments()" viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2"
          fill="none" stroke-linecap="round" stroke-linejoin="round" class="css-i6dzq1">
          <path d="M12 20h9"></path>
          <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
        </svg>
        {{ this.imageData.comments }}
      </span>
    </div>
  </div>


</template>


<style>

</style>

