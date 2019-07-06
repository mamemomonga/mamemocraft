function Application() {}
Application.prototype={
	run: function() {
		const t=this
		t.monitorTimer=60*30
		t.monitorTimerRemain=t.monitorTimer
		$("#power").on("click",()=>{ t.poweron() })
		t.interval=setInterval(()=>{ t.state() },1000)
	},
	state: function() {
		const t=this
		$.get("/api/state",(r)=>{
			// console.log(r)
			switch(r.state) {
				case 0: // unknown
					$("#power").hide()
					t.stateMark("stop")
					break
				case 1: // stop
					$("#power").show()
					t.stateMark("stop")
					break
				case 2: // loading
					$("#power").hide()
					t.stateMark("load")
					break
				case 3: // runnning
					$("#power").hide()
					t.stateMark("run")
					break
			}
			$("#status_message").text(r.message)
		},"json");
		t.monitorTimerRemain--;
		if(t.monitorTimerRemain <= 0) {
			t.stopController()
		}
	},
	poweron: function() {
		const t=this
		t.monitorTimerRemain=t.monitorTimer;
		console.log("POWERON")
		$("#power").hide()
		$.get("/api/poweron")
	},
	stopController: function() {
		const t=this;
		console.log("stopController")
		clearInterval(t.interval)
		setTimeout(()=> {
			$("#power").hide()
			t.stateMark("none")
			$("#status_message").text("このページをリロードしてください")
		},1000)
	},
	stateMark: function(mark) {
		switch(mark) {
			case "stop":
				$("#status_stop").show()
				$("#status_run").hide()
				$("#status_load").hide()
				break
			case "run":
				$("#status_stop").hide()
				$("#status_run").show()
				$("#status_load").hide()
				break
			case "load":
				$("#status_stop").hide()
				$("#status_run").hide()
				$("#status_load").show()
				break
			case "none":
				$("#status_stop").hide()
				$("#status_run").hide()
				$("#status_load").hide()
				break

		}
	}

}

$(function() { new Application().run() })

