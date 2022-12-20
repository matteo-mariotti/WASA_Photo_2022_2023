<script>

import backend from "@/services/axios";

export default {

  data: function () {
    return {
      file: null,
    }
  },
  methods: {
    handleFileUpload() {
      this.file = this.$refs.file.files[0];
      console.log(this.file)
    },
    async submitFile() {
      try {
        let formData = new FormData()
        formData.append("file", this.file)

        await backend.post(`/users/${sessionStorage.getItem("username")}/photos`,formData)

      }catch(error){
        console.log(error)
      }
    },
  },
}


</script>


<template>
  <div class="container">
    <div class="text-center">
      <h5>Choose a photo to upload</h5>
      <hr/>
      <form enctype="multipart/form-data" @submit.prevent="submitFile">
      <input type="file" class="form-control" ref="file" @change="handleFileUpload"/>
      <br>
      <button class="btn btn-primary">Upload</button>
      </form>
    </div>
  </div>
</template>

