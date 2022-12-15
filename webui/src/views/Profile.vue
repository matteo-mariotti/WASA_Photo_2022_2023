<script>
import backend from "@/services/axios";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SuccessMsg from "@/components/SuccessMsg.vue";
import PhotoBox from "@/components/PhotoBox.vue";

export default {
  components: {PhotoBox, SuccessMsg, ErrorMsg},
  data: function () {
    return {
      errormsg: null,
      ready: false,
      userData: null,
      successmsg: null,
      userToken: this.$route.params.userID
    }
  },
  methods: {
    async loadContent() {
      try {
        this.ready = false
        let response = await backend.get(`/users/${this.$route.params.userID}`); // TODO Sitemare i parametri
        this.handleResponse(response);
        this.ready = true

      } catch (error) {
        this.handleError(error);
      }
    },
    handleResponse(response) {
      this.errormsg = null
      this.userData = response.data
      console.log(response)
    },
    handleError(error) {

      // Print the error from the server (for debugging)
      this.errormsg = error.response.data.message

      // Log the error on console
      console.log(error)
    }
  },
  mounted() {
    this.loadContent();
  }
}
</script>

<template>


  <div class="container mt-5">


    <div class="row">
      <div class="col-3"></div>
      <div class="col-6">

        <h1 class="h2" style="text-align: center" v-if="ready">{{ userData.username }}</h1>

        <div class="row" style="text-align: center">
          <!--TODO Rendere clickable le label qui sotto per aprire la lista dei followers e dei seguiti -->

          <div class="col-4">
            <label v-if="ready">Photo:{{ userData.photoNumber }}</label>
          </div>
          <div class="col-4">
            <label v-if="ready"> Follower: {{ userData.follower }}</label>
          </div>
          <div class="col-4">
            <label v-if="ready"> Following: {{ userData.following }}</label>
          </div>

        </div>


        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>
      </div>

    </div>

    <div class="row" style="background-color: red; align-content: center">
      <div class="col-3"></div>
      <div class="col-6 align-items-center d-grid gap-5">

        <div v-if="ready" v-for="image in userData.photos" id="photoList">
          <PhotoBox v-bind:image-data="image.id"></PhotoBox>
        </div>

        <br>
      </div>


      <div class="col-3"></div>


    </div>



  </div>
</template>

<style>
</style>
