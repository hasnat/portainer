angular.module('portainer.app')
.controller('CommandController', ['$scope', '$state', '$transition$', '$filter', 'CommandService', 'Notifications',
function ($scope, $state, $transition$, $filter, CommandService, Notifications) {

  if (!$scope.applicationState.application.commandManagement) {
    $state.go('portainer.commands');
  }

  $scope.state = {
    actionInProgress: false
  };

  $scope.updateEndpoint = function() {
    var command = $scope.command;
    var commandParams = {
      name: command.Name,
      image: command.Image,
      command: command.Command
    };

    $scope.state.actionInProgress = true;
    CommandService.updateEndpoint(command.Id, commandParams)
    .then(function success(data) {
      Notifications.success('Command updated', $scope.endpoint.Name);
      $state.go('portainer.commands');
    }, function error(err) {
      Notifications.error('Failure', err, 'Unable to update command');
      $scope.state.actionInProgress = false;
    }, function update(evt) {
      if (evt.upload) {
        $scope.state.uploadInProgress = evt.upload;
      }
    });
  };

  function initView() {
    CommandService.endpoint($transition$.params().id)
    .then(function success(endpoint) {
      $scope.endpoint = endpoint;
    })
    .catch(function error(err) {
      Notifications.error('Failure', err, 'Unable to retrieve command details');
    });
  }

  initView();
}]);
