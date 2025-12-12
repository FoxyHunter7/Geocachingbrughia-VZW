<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('geocaches', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->string('state');
            $table->string('title');
            $table->string('type');
            $table->string('geolink');
            $table->integer('difficulty');
            $table->integer('terrain');
            $table->dateTime('placed_on');

            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('geocaches');
    }
};
