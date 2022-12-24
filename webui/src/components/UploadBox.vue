<script>

import backend from "@/services/axios";
import SuccessMsg from "@/components/SuccessMsg.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";

export default {
  components: {ErrorMsg, SuccessMsg},
  data: function () {
    return {
      file: null,
      successmsg: null,
      errormsg: null,
    }
  },
  methods: {
    handleFileUpload() {
      this.file = this.$refs.file.files[0];
      if (!this.fileValidation()){
        this.errormsg = "Invalid file type, allowed types are .jpg,.jpeg,.png,.webp"
        return
      }
      document.getElementById("uploadBtn").disabled = false;
    },
    reload() {
      location.reload()
    },
    async submitFile() {
      try {

        let formData = new FormData()
        formData.append("file", this.file)

        await backend.post(`/users/${sessionStorage.getItem("username")}/photos`, formData)
        this.errormsg = null
        this.successmsg = "Photo successfully uploaded"
        setTimeout(this.reload, 1000)
      } catch (error) {
        this.errormsg = "Error while uploading photo"
        console.log(error)
      }
    },
    fileValidation() {
      let fileInput = document.getElementById('file');

      let filePath = fileInput.value;

      // Allowing file type
      let allowedExtensions = /(\.jpg|\.jpeg|\.png|\.webp)$/i;

      if (!allowedExtensions.exec(filePath)) {
        fileInput.value = '';
        return false;
      }
      return true;
    },
  },
}


</script>


<template>
  <div class="container">

    <div class="text-center">
      <h5>Choose a photo to upload</h5>
      <hr>
      <form enctype="multipart/form-data" @submit.prevent="submitFile">
        <input type="file" id="file" class="form-control" ref="file" @change="handleFileUpload" accept=".jpg,.jpeg,.png,.webp" />
        <br>
        <button id="uploadBtn" class="btn btn-primary" disabled>Upload</button>
      </form>
      <SuccessMsg :msg="successmsg" v-if="successmsg" class="mt-2"></SuccessMsg>
      <ErrorMsg :msg="errormsg" v-if="errormsg" class="mt-2"></ErrorMsg>
    </div>
  </div>
</template>

