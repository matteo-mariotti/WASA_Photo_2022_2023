<script>
import backend from "@/services/axios";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SuccessMsg from "@/components/SuccessMsg.vue";
import Card from "@/components/Card.vue";

export default {
  components: {Card, SuccessMsg, ErrorMsg},
  data: function () {
    return {
      errormsg: null,
      ready: false,
      userData: null,
      successmsg: null,
      userToken: this.$route.params.userID,
      currPage: 0
    }
  },
  methods: {
    async loadContent() {
      try {
        this.ready = false
        let response = await backend.get(`/users/${sessionStorage.getItem("token")}`); // TODO Sitemare i parametri
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
    },
    loadPhotos(offset){
      // TODO Fare la chiamata per ottenere le foto a partire da un certo offset

    },
  },
  mounted() {
    this.loadContent();
  }
}
</script>

<template>

  <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

  <div class="justify-content-center align-items-center">

    <!-- Username of the user -->
    <div class="h2 d-flex justify-content-evenly mt-4" v-if="ready" style="font-size:3.5vw;">

      {{userData.username}}
      <button class="btn btn-outline-primary btn-sm"> Change my username</button>

    </div>
    <!-- Profile stats -->
    <div class="d-flex justify-content-around p-2" style="font-size: 2vw;background-color: rgba(137,137,137,0.2);border-radius: 2vw;">
      <div v-if="ready">Photo: {{ userData.photoNumber }}</div>
      <div v-if="ready">Follower: {{ userData.follower }}</div>
      <div v-if="ready">Following: {{ userData.following }}</div>
    </div>

    <!-- Images -->
    <div class="d-grid justify-content-center" id="imageList">
      <Card v-bind:imageData="image" v-if="ready" v-for="image in userData.photos"></Card>
    </div>
    <div class="d-grid justify-content-center">
      <button @click="loadPhotos(this.currentOffset)" type="button" class="btn btn-outline-dark m-3">Load More</button>
    </div>


  </div>

</template>

<style>
</style>
