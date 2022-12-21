<script>
import ErrorMsg from "@/components/ErrorMsg.vue";
import SuccessMsg from "@/components/SuccessMsg.vue";
import UploadBox from "@/components/UploadBox.vue";

export default {
  components: {UploadBox, SuccessMsg, ErrorMsg},
  data: function () {
    return {
      errormsg: null,
      successmsg: null
    }
  },
  methods: {},
}
</script>

<template>


  <div class="container mt-5">
    <div class="row">
      <div class="col-sm-3"></div>
      <div class="col-sm-6">

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

        <UploadBox></UploadBox>

      </div>

      <div class="col-sm-3"></div>

    </div>

  </div>
</template>

<style>
::placeholder {
  text-align: center;
}
</style>
