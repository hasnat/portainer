angular.module('portainer.app')
.factory('CommandService', ['$q', 'Commands',
function CommandServiceFactory($q, Commands) {
  'use strict';
  var service = {};

  service.command = function(id) {
    return Commands.get({id}).$promise;
  };

  service.commands = function() {
    return Commands.query({}).$promise;
  };

  service.updateCommand = function(id, command) {
    return Commands.update({id}, command).$promise;
  };

  service.deleteCommand = function(id) {
    return Commands.remove({id}).$promise;
  };

  service.createCommand = function(name, image, command) {
    return Endpoints.create({}, {name, image, command}).$promise;
  };

  return service;
}]);
