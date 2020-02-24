// @ts-nocheck
var status = rs.status();

if (status.codeName === 'NotYetInitialized') {
  print('Replica Set not yet initialized. Initializing...');

  rs.initiate({
    _id: 'rs0',
    members: [
      {
        _id: 0,
        host: 'mdb:27017',
      },
    ],
  });

  print('Replica Set successfully initialized. Replica Set Ready.');
} else {
  print('Replica Set already initialized. Replica Set Ready.');
}
