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
      // datetime: this.imageData. TODO Aggiungere la data
    }
  },
  methods: {
    async loadImage() {
      try {
        this.ready = false
        let response = await backend.get(`/photos/${this.imageData}`, {responseType: 'arraybuffer'});

        console.log("Request executed");
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
      let img = document.getElementById(this.imageData)
      let blob = new Blob([response.data])
      img.src = URL.createObjectURL(blob)

      console.log("Blob created successfully")

    },
    handleError(error) {

      // Print the error from the server (for debugging)
      this.errormsg = error.response.data.message

      // Log the error on console
      console.log(error)
    },
    test() {
      alert("Like test")
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
    <img v-bind:id="this.imageData" class="rounded-5 img-fluid" alt="User image">

    <div class="row">
      <div style="background-color:#ffea81" class="col-4">
        <div class="w-25" style="background-color: #86b7ff">
          Likes
          <svg @click="test()" xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="bi bi-star-center"
               viewBox="1 1 16 16">
            <path
                d="M2.866 14.85c-.078.444.36.791.746.593l4.39-2.256 4.389 2.256c.386.198.824-.149.746-.592l-.83-4.73 3.522-3.356c.33-.314.16-.888-.282-.95l-4.898-.696L8.465.792a.513.513 0 0 0-.927 0L5.354 5.12l-4.898.696c-.441.062-.612.636-.283.95l3.523 3.356-.83 4.73zm4.905-2.767-3.686 1.894.694-3.957a.565.565 0 0 0-.163-.505L1.71 6.745l4.052-.576a.525.525 0 0 0 .393-.288L8 2.223l1.847 3.658a.525.525 0 0 0 .393.288l4.052.575-2.906 2.77a.565.565 0 0 0-.163.506l.694 3.957-3.686-1.894a.503.503 0 0 0-.461 0z"/>
          </svg>
        </div>
      </div>
      <div class="col-4" style="background-color:#86b7fe ">Mezzo</div>
      <div style="background-color:#0dcaf0 " class="col-4">
        Comments


      </div>
    </div>
  </div>
</template>

<style></style>
