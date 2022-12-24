<script>

import backend from "@/services/axios";

export default {
  props: ["comment", "photo"],
  data: function () {
    return {
      isOwned: this.comment.userID === sessionStorage.getItem("username"),
      commentNumber: this.photo.comments,
    }
  },
  methods: {
    async deleteComment(){
      try {
        await backend.delete(`/users/${this.photo.photoOwner}/photos/${this.photo.id}/comments/${this.comment.commentID}`)
        this.commentNumber = this.commentNumber-1
        this.$emit('update', this.comment)
        this.$emit('updateCount', this.comment)
      }catch(error){
        console.log(error)
      }
    },
  },
}

</script>


<template>

  <div class="card m-2">
    <div class="card-body">
        <div class="d-flex">
          <h4 class="card-title">{{this.comment.userID}}</h4>
          <button @click="deleteComment()" v-if="isOwned" class="btn btn-outline-danger btn-sm ms-auto"> Delete comment
          </button>
        </div>
      <p class="card-text">{{this.comment.text}}</p>
    </div>
  </div>

</template>

<style>
</style>