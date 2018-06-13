angular.module('portainer.app')
.factory('CommandProvider', ['LocalStorage', function CommandProviderFactory(LocalStorage) {
  'use strict';
  var service = {};
  var command = {};

  service.initialize = function() {

  };

  service.clean = function() {
    command = {};
  };


  return service;
}]);
