{
  // Describe containers
  containers: [ // this actually belongs to the dictionary below
    {
      /*
      we need to find the right value
      this will be set at runtime
      */
      foo: ${magic}
      memory: 256.4 // this is an inline comment
      required: true
      path: "${basepath}/repo"
      image: "golang:${tag}-node"
      value: /* this is some crazy comment */ "string\nwith\t\"escapes\\asdasd"
      portMappings: [{
        "hostPort": 8080,
        "containerPort": null, /* TODO what should this be */
        `protocol`: `tcp`,
        }]
      // This is a foot comment
    } // this is actually inline with the dictionary
    // this is a foot comment too
  ] // this is actually inline with the list
  description: `this is a thing
that does stuff
over multiple lines!` // this is inline with the raw string
  emptylist: [
  ]
  emptymap: {

      }
  "list-with-comments": [ /* this is part of the list */
    // where does this go?
  ]
  "map-with-comments": {
    // nothing to see here
  }
}, // trailing comma on root because I can
// foot comment on root