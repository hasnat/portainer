angular.module('portainer.docker')
.controller('LogViewerController', ['clipboard',
function (clipboard) {
  var ctrl = this;

  this.state = {
    copySupported: clipboard.supported,
    logCollection: true,
    autoScroll: true,
    wrapLines: true,
    search: '',
    filteredLogs: [],
    selectedLines: []
  };

  this.copy = function() {
    clipboard.copyText(this.state.filteredLogs);
    $('#refreshRateChange').show();
    $('#refreshRateChange').fadeOut(2000);
  };

  this.copySelection = function() {
    clipboard.copyText(this.state.selectedLines);
    $('#refreshRateChange').show();
    $('#refreshRateChange').fadeOut(2000);
  };

  this.clearSelection = function() {
    this.state.selectedLines = [];
  };

  this.selectLine = function(line) {
    var idx = this.state.selectedLines.indexOf(line);
    if (idx === -1) {
      this.state.selectedLines.push(line);
    } else {
      this.state.selectedLines.splice(idx, 1);
    }
  };

  this.getLogViewerCSSClass = function() {
    return 'log_viewer' + (this.state.wrapLines ? ' wrap_lines' : '');
  };
}]);
