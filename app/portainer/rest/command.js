angular.module('portainer.app')
.factory('Commands', ['$resource', 'API_ENDPOINT_COMMANDS', function CommandsFactory($resource, API_ENDPOINT_COMMANDS) {
  'use strict';
  return $resource(API_ENDPOINT_COMMANDS + '/:id/:action', {}, {
    create: { method: 'POST', ignoreLoadingBar: true },
    query: { method: 'GET', isArray: true },
    get: { method: 'GET', params: { id: '@id' } },
    update: { method: 'PUT', params: { id: '@id' } },
    remove: { method: 'DELETE', params: { id: '@id'} }
  });
}]);
