# WHY
This package is meant to be a lightweight package that both the gopherjs-based UI and the backend can import to share communication constants.

It is a separate, stand-alone, lightweight package to avoid compiling big parts of the frontend in the backend.

One more reason to keep these constants all together is to make it easier to build alternative UIs for wapty and to understand which functionalities are exposed by the websocket api.
